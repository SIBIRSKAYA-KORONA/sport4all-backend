package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type TournamentUseCaseImpl struct {
	userRepo       repositories.UserRepository
	tournamentRepo repositories.TournamentRepository
	meetingRepo    repositories.MeetingRepository
}

func CreateTournamentUseCase(userRepo repositories.UserRepository, tournamentRepo repositories.TournamentRepository,
	meetingRepo repositories.MeetingRepository) usecases.TournamentUseCase {
	return &TournamentUseCaseImpl{
		userRepo:       userRepo,
		tournamentRepo: tournamentRepo,
		meetingRepo:    meetingRepo,
	}
}

func (tournamentUseCase *TournamentUseCaseImpl) Create(ownerId uint, tournament *models.Tournament) error {
	if _, err := tournamentUseCase.userRepo.GetByID(ownerId); err != nil { // TODO: move it to mv
		logger.Error(err)
		return err
	}

	tournament.OwnerId = ownerId // TODO: init in handler
	if err := tournamentUseCase.tournamentRepo.Create(tournament); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) GetByID(tid uint) (*models.Tournament, error) {
	tournament, err := tournamentUseCase.tournamentRepo.GetByID(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tournament, nil
}

func (tournamentUseCase *TournamentUseCaseImpl) AddTeam(tournamentId uint, teamId uint) error {
	if err := tournamentUseCase.tournamentRepo.AddTeam(tournamentId, teamId); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) GetAllTeams(tournamentId uint) (*models.Teams, error) {
	teams, err := tournamentUseCase.tournamentRepo.GetAllTeams(tournamentId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return teams, nil
}

func (tournamentUseCase *TournamentUseCaseImpl) GenerateMeetings(tournamentId uint, genType uint) error {
	switch genType {
	case usecases.Olympic:
		return tournamentUseCase.generateOlympicMeetings(tournamentId)
	case usecases.Circular:
		return tournamentUseCase.generateCircularMeetings(tournamentId)
	default:
		return errors.ErrInternal // TODO create error for this event
	}
}

func (tournamentUseCase *TournamentUseCaseImpl) generateOlympicMeetings(tournamentId uint) error {
	teams, err := tournamentUseCase.GetAllTeams(tournamentId)
	if err != nil {
		logger.Error(err)
		return err
	}

	if len(*teams)%2 != 0 { // TODO: add valid param
		logger.Error("invalid tournament size")
		return errors.ErrInternal // TODO create error for this event
	}

	for i := 0; i < len(*teams); i += 2 {
		err = tournamentUseCase.meetingRepo.Create(
			&models.Meeting{
				Status:       usecases.New,
				TournamentId: tournamentId,
				Teams:        []models.Team{(*teams)[i], (*teams)[i+1]},
			})

		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) generateCircularMeetings(tournamentId uint) error {
	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) GetAllMeetings(tournamentId uint) (*models.Meetings, error) {
	meetings, err := tournamentUseCase.tournamentRepo.GetAllMeetings(tournamentId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return meetings, nil
}
