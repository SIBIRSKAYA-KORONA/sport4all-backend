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
	teamRepo       repositories.TeamRepository
	meetingRepo    repositories.MeetingRepository
}

func CreateTournamentUseCase(userRepo repositories.UserRepository, tournamentRepo repositories.TournamentRepository,
	teamRepo repositories.TeamRepository, meetingRepo repositories.MeetingRepository) usecases.TournamentUseCase {
	return &TournamentUseCaseImpl{
		userRepo:       userRepo,
		tournamentRepo: tournamentRepo,
		teamRepo:       teamRepo,
		meetingRepo:    meetingRepo,
	}
}

func (tournamentUseCase *TournamentUseCaseImpl) Create(tournament *models.Tournament) error {
	if _, err := tournamentUseCase.userRepo.GetByID(tournament.OwnerId); err != nil { // TODO: move it to mv
		logger.Error(err)
		return err
	}

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

func (tournamentUseCase *TournamentUseCaseImpl) GetTournamentByUser(uid uint) (*models.UserTournament, error) {
	userTournament := new(models.UserTournament)
	tournament, err := tournamentUseCase.tournamentRepo.GetTournamentByUser(uid)
	if err != nil {
		logger.Error(err)
	}
	userTournament.Owner = *tournament

	teams, err := tournamentUseCase.teamRepo.GetTeamsByUser(uid, models.Player)
	if err != nil {
		logger.Error(err)
	}

	for _, team := range *teams {
		tournament, err = tournamentUseCase.teamRepo.GetAllTournaments(team.ID)
		if err != nil {
			logger.Error(err)
		}
		userTournament.TeamMember = append(userTournament.TeamMember, *tournament...)
	}

	teams, err = tournamentUseCase.teamRepo.GetTeamsByUser(uid, models.Owner)
	if err != nil {
		logger.Error(err)
	}

	for _, team := range *teams {
		tournament, err = tournamentUseCase.teamRepo.GetAllTournaments(team.ID)
		if err != nil {
			logger.Error(err)
		}
		userTournament.Owner = append(userTournament.Owner, *tournament...)
	}

	return userTournament, nil
}

func (tournamentUseCase *TournamentUseCaseImpl) Update(tournament *models.Tournament) error {
	// TODO: validate tournament.Status
	if tournament.Status == models.InProgressEvent {
		// TODO: may be sent tournament system from front
		oldTournament, err := tournamentUseCase.GetByID(tournament.ID)
		if err != nil {
			return err
		}
		if err = tournamentUseCase.generateMesh(oldTournament.ID, oldTournament.System); err != nil {
			logger.Error(err)
			return err
		}
	}

	if err := tournamentUseCase.tournamentRepo.Update(tournament); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) AddTeam(tournamentId uint, teamId uint) error {
	tournament, err := tournamentUseCase.GetByID(tournamentId)
	if err != nil {
		return err
	}
	if tournament.Status > models.RegistrationEvent {
		return errors.ErrInternal // TODO: create error for this event
	}

	if err = tournamentUseCase.tournamentRepo.AddTeam(tournamentId, teamId); err != nil {
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

func (tournamentUseCase *TournamentUseCaseImpl) GetAllMeetings(tournamentId uint) (*models.Meetings, error) {
	meetings, err := tournamentUseCase.tournamentRepo.GetAllMeetings(tournamentId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return meetings, nil
}

func (tournamentUseCase *TournamentUseCaseImpl) generateMesh(tournamentId uint, genType models.TournamentSystem) error {
	switch genType {
	case models.OlympicSystem:
		return tournamentUseCase.generateOlympicMesh(tournamentId)
	case models.CircularSystem:
		return tournamentUseCase.generateCircularMesh(tournamentId)
	default:
		return errors.ErrInternal // TODO create error for this event
	}
}

func (tournamentUseCase *TournamentUseCaseImpl) generateOlympicMesh(tournamentId uint) error {
	teams, err := tournamentUseCase.GetAllTeams(tournamentId)
	if err != nil {
		logger.Error(err)
		return err
	}

	size := len(*teams) / 2
	for i := 0; i < size; i++ {
		// TODO: make save batch
		// TODO: create balanced bin tree
		meeting := &models.Meeting{
			Status:       models.NotStartedEvent,
			Round:        0,
			Group:        0,
			TournamentId: tournamentId,
		}
		if err = tournamentUseCase.meetingRepo.Create(meeting); err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) generateCircularMesh(tournamentId uint) error {
	panic("not implement method generateCircularMesh")
}
