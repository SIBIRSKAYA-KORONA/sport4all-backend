package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type TeamUseCaseImpl struct {
	teamRepo repositories.TeamRepository
	userRepo repositories.UserRepository
}

func CreateTeamUseCase(teamRepo repositories.TeamRepository, userRepo repositories.UserRepository) usecases.TeamUseCase {
	return &TeamUseCaseImpl{teamRepo: teamRepo, userRepo: userRepo}
}

func (teamUseCase *TeamUseCaseImpl) Create(ownerId uint, team *models.Team) error {
	// TODO: кажется этот запрос не нужен или его надо вынести в mv (Антон)
	if _, err := teamUseCase.userRepo.GetByID(ownerId); err != nil {
		logger.Error(err)
		return err
	}

	team.OwnerId = ownerId

	if err := teamUseCase.teamRepo.Create(team); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (teamUseCase *TeamUseCaseImpl) GetByID(tid uint) (*models.Team, error) {
	team, err := teamUseCase.teamRepo.GetByID(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return team, nil
}

func (teamUseCase *TeamUseCaseImpl) CheckUserForRole(teamID uint, userID uint, role models.Role) (bool, error) {
	switch role {
	case models.Owner:
		return teamUseCase.teamRepo.IsTeamOwner(teamID, userID)
	case models.Player:
		return teamUseCase.teamRepo.IsTeamPlayer(teamID, userID)
	default:
		return false, errors.ErrTeamBadRole
	}
}

func (teamUseCase *TeamUseCaseImpl) GetTeamsByUser(uid uint, role models.Role) (*models.Teams, error) {
	teams, err := teamUseCase.teamRepo.GetTeamsByUser(uid, role)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return teams, nil
}

func (teamUseCase *TeamUseCaseImpl) GetAllTournaments(tid uint) (*models.Tournaments, error) {
	tournaments, err := teamUseCase.teamRepo.GetAllTournaments(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return tournaments, nil
}

func (teamUseCase *TeamUseCaseImpl) GetTeamsByNamePart(namePart string, limit uint) (*models.Teams, error) {
	users, err := teamUseCase.teamRepo.GetTeamsByNamePart(namePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (teamUseCase *TeamUseCaseImpl) GetUsersForInvite(tid uint, nicknamePart string, limit uint) (*models.Users, error) {
	users, err := teamUseCase.teamRepo.GetUsersForInvite(tid, nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (teamUseCase *TeamUseCaseImpl) InviteMember(tid uint, uid uint, role models.Role) error {
	user, err := teamUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err = teamUseCase.teamRepo.InviteMember(tid, user, role); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (teamUseCase *TeamUseCaseImpl) DeleteMember(tid uint, uid uint) error {
	if _, err := teamUseCase.userRepo.GetByID(uid); err != nil {
		logger.Error(err)
		return err
	}

	if err := teamUseCase.teamRepo.DeleteMember(tid, uid); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (teamUseCase *TeamUseCaseImpl) GetTeamStats(tid uint) ([]models.Stats, error) {
	stats, err := teamUseCase.teamRepo.GetTeamStats(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return stats, nil
}
