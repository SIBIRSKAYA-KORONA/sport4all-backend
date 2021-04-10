package impl

import (
	"math"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type TournamentUseCaseImpl struct {
	userRepo                 repositories.UserRepository
	tournamentRepo           repositories.TournamentRepository
	teamRepo                 repositories.TeamRepository
	meetingRepo              repositories.MeetingRepository
	tournamentSystemsCreator map[string]func(uint) error
}

func CreateTournamentUseCase(userRepo repositories.UserRepository, tournamentRepo repositories.TournamentRepository,
	teamRepo repositories.TeamRepository, meetingRepo repositories.MeetingRepository) usecases.TournamentUseCase {

	impl := &TournamentUseCaseImpl{
		userRepo:       userRepo,
		tournamentRepo: tournamentRepo,
		teamRepo:       teamRepo,
		meetingRepo:    meetingRepo,
	}

	impl.tournamentSystemsCreator = map[string]func(uint) error{
		models.OlympicSystem:  impl.generateOlympicMesh,
		models.CircularSystem: impl.generateCircularMesh,
	}

	return impl
}

func (tournamentUseCase *TournamentUseCaseImpl) Create(tournament *models.Tournament) error {
	if _, ok := tournamentUseCase.tournamentSystemsCreator[tournament.System]; !ok {
		return errors.ErrTournamentSystemNotAcceptable
	}

	// TODO: move it to mv
	if _, err := tournamentUseCase.userRepo.GetByID(tournament.OwnerId); err != nil {
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
	tournament, err := tournamentUseCase.tournamentRepo.GetTournamentByUserOwner(uid)
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

func (tournamentUseCase *TournamentUseCaseImpl) CheckUserForTournamentRole(tournamentId uint,
	uid uint, role models.TournamentRole) (bool, error) {

	switch role {
	case models.TournamentOrganizer:
		return tournamentUseCase.tournamentRepo.IsTournamentOrganizer(tournamentId, uid)
	case models.TournamentPlayer:
		return tournamentUseCase.tournamentRepo.IsTournamentPlayer(tournamentId, uid)
	default:
		return false, errors.ErrTournamentBadRole
	}
}

func (tournamentUseCase *TournamentUseCaseImpl) Update(tournament *models.Tournament) error {
	// TODO: move it to mv (Антон)
	oldTournament, err := tournamentUseCase.GetByID(tournament.ID)
	if err != nil {
		return err
	}
	if oldTournament.Status > models.NotStartedEvent && tournament.Status <= oldTournament.Status {
		return errors.ErrTournamentStatusNotAcceptable
	}

	switch tournament.Status {
	case models.UnknownEvent, models.NotStartedEvent:
		if err = tournamentUseCase.tournamentRepo.Update(tournament); err != nil {
			logger.Error(err)
			return err
		}
	case models.InProgressEvent:
		generateMesh, ok := tournamentUseCase.tournamentSystemsCreator[oldTournament.System]
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

func (tournamentUseCase *TournamentUseCaseImpl) IsTeamInTournament(tournamentId uint, teamId uint) (bool, error) {
	result, err := tournamentUseCase.tournamentRepo.IsTeamInTournament(tournamentId, teamId)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return result, err
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

	for i := range *meetings {
		if (*meetings)[i].Status > models.RegistrationEvent {
			stat, err := tournamentUseCase.meetingRepo.GetMeetingTeamStat((*meetings)[i].ID)
			if err != nil {
				logger.Warn(err)
			}
			(*meetings)[i].Stats = *stat
		}
	}

	return meetings, nil
}

func generateOlympicMeshImpl(root *models.Meeting, deep int) {
	root.Status = models.RegistrationEvent
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
		return err
	}

	root := &models.Meeting{TournamentId: tournamentId, NextMeetingID: nil}
	generateOlympicMeshImpl(root, int(math.Ceil(math.Log2(float64(len(*teams))))))

	if err = tournamentUseCase.meetingRepo.Create(root); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) calcCircularRound(teamNum1, teamNum2, maxRounds int) uint {
	round := 0
	sum := teamNum1 + teamNum2
	if maxRounds == teamNum1 || maxRounds == teamNum2 {
		sum = teamNum1
		if teamNum1 > teamNum2 {
			sum = teamNum2
		}
		sum *= 2
	}

	if sum <= maxRounds {
		round = sum - 1
	} else {
		round = sum - maxRounds
	}

	return uint(round)
}

func (tournamentUseCase *TournamentUseCaseImpl) generateCircularMesh(tournamentId uint) error {
	teams, err := tournamentUseCase.GetAllTeams(tournamentId)
	if err != nil {
		return err
	}

	numTeams := len(*teams)
	numRound := numTeams
	if numRound%2 != 0 {
		numRound++
	}
	var meetings []models.Meeting
	for i := 0; i < numTeams; i++ {
		for j := 0; j < i; j++ {
			meetings = append(meetings, models.Meeting{
				Status:       models.RegistrationEvent,
				Group:        0,
				Round:        tournamentUseCase.calcCircularRound(i+1, j+1, numRound),
				TournamentId: tournamentId,
				Teams:        []models.Team{(*teams)[i], (*teams)[j]},
			})
		}
	}
	logger.Debug("generate circular mesh tournamentId: ", tournamentId, " , with ", len(meetings), " meetings")

	for idx := range meetings {
		if err = tournamentUseCase.meetingRepo.Create(&meetings[idx]); err != nil {
			logger.Error(err)
		}
	}

	/*
		 // TODO: поч не работает Create Batch
			if err = tournamentUseCase.meetingRepo.CreateBatch(&meetings); err != nil {
				logger.Error(err)
				return err
			}
	*/

	return nil
}

func (tournamentUseCase *TournamentUseCaseImpl) GetTournamentForFeeds(offset, maxTournament uint) (*[]models.Tournament, error) {
	result, err := tournamentUseCase.tournamentRepo.GetTournamentForFeeds(offset, maxTournament)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return result, nil
}
