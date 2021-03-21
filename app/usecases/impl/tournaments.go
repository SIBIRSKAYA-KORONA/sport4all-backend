package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type TournamentUseCaseImpl struct {
	userRepo          repositories.UserRepository
	tournamentRepo    repositories.TournamentRepository
	teamRepo          repositories.TeamRepository
	meetingRepo       repositories.MeetingRepository
	tournamentSystems map[string]func(uint) error
}

func CreateTournamentUseCase(userRepo repositories.UserRepository, tournamentRepo repositories.TournamentRepository,
	teamRepo repositories.TeamRepository, meetingRepo repositories.MeetingRepository) usecases.TournamentUseCase {

	impl := &TournamentUseCaseImpl{
		userRepo:       userRepo,
		tournamentRepo: tournamentRepo,
		teamRepo:       teamRepo,
		meetingRepo:    meetingRepo,
	}

	impl.tournamentSystems = map[string]func(uint) error{
		"olympic":  impl.generateOlympicMesh,
		"circular": impl.generateCircularMesh,
	}

	return impl
}

func (tournamentUseCase *TournamentUseCaseImpl) Create(tournament *models.Tournament) error {
	if _, ok := tournamentUseCase.tournamentSystems[tournament.System]; !ok {
		return errors.ErrTournamentSystemNotAcceptable
	}

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
	// TODO: move it to mv (Антон)
	oldTournament, err := tournamentUseCase.GetByID(tournament.ID)
	if err != nil {
		return err
	}
	if oldTournament.Status > models.NotStartedEvent && tournament.Status < oldTournament.Status {
		return errors.ErrTournamentStatusNotAcceptable
	}

	switch tournament.Status {
	case models.UnknownEvent, models.NotStartedEvent:
		if err = tournamentUseCase.tournamentRepo.Update(tournament); err != nil {
			logger.Error(err)
			return err
		}
	case models.InProgressEvent:
		generateMesh, ok := tournamentUseCase.tournamentSystems[oldTournament.System]
		if !ok {
			return errors.ErrTournamentSystemNotAcceptable
		}

		if err = generateMesh(oldTournament.ID); err != nil {
			return err
		}
		fallthrough
	case models.RegistrationEvent, models.FinishedEvent:
		if err = tournamentUseCase.tournamentRepo.
			Update(&models.Tournament{ID: tournament.ID, Status: tournament.Status}); err != nil {
			logger.Error(err)
			return err
		}
	default:
		return errors.ErrTournamentStatusNotAcceptable
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) AddTeam(tournamentId uint, teamId uint) error {
	// TODO: move it to mv (Антон)
	tournament, err := tournamentUseCase.GetByID(tournamentId)
	if err != nil {
		return err
	}
	if tournament.Status != models.RegistrationEvent {
		return errors.ErrTournamentStatusNotAcceptable
	}

	if err = tournamentUseCase.tournamentRepo.AddTeam(tournamentId, teamId); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) RemoveTeam(tournamentId uint, teamId uint) error {
	// TODO: move it to mv (Антон)
	tournament, err := tournamentUseCase.GetByID(tournamentId)
	if err != nil {
		return err
	}
	if tournament.Status != models.RegistrationEvent {
		return errors.ErrTournamentStatusNotAcceptable
	}

	if err = tournamentUseCase.tournamentRepo.RemoveTeam(tournamentId, teamId); err != nil {
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

func generateOlympicMeshImpl(root *models.Meeting, deep int) {
	root.Status = models.NotStartedEvent
	root.Round = uint(deep)
	root.Group = 0

	deep--
	if deep <= 0 {
		return
	}

	root.PrevMeetings = make([]models.Meeting, 2)
	for idx := range root.PrevMeetings {
		root.PrevMeetings[idx].NextMeetingID = &root.ID
		root.PrevMeetings[idx].TournamentId = root.TournamentId
		generateOlympicMeshImpl(&root.PrevMeetings[idx], deep)
	}
}

func (tournamentUseCase *TournamentUseCaseImpl) generateOlympicMesh(tournamentId uint) error {
	teams, err := tournamentUseCase.GetAllTeams(tournamentId)
	if err != nil {
		logger.Error(err)
		return err
	}

	root := &models.Meeting{TournamentId: tournamentId /*, NextMeeting: nil*/}
	numTeams := len(*teams)
	if numTeams%2 != 0 {
		numTeams++
	}
	generateOlympicMeshImpl(root, numTeams/2)

	if err = tournamentUseCase.meetingRepo.Create(root); err != nil {
		return err
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) generateCircularMesh(tournamentId uint) error {
	// TODO: напиши меня (Тим)
	panic("not implement method generateCircularMesh")
}
