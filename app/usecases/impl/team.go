package impl

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type TeamUseCaseImpl struct {
	teamRepo repositories.TeamRepository
	userRepo repositories.UserRepository
}

func CreateTeamUseCase(teamRepo repositories.TeamRepository, userRepo repositories.UserRepository) usecases.TeamUseCase {
	return &TeamUseCaseImpl{teamRepo: teamRepo, userRepo: userRepo}
}

func (teamUseCase *TeamUseCaseImpl) Create(ownerId uint, team *models.Team) error {
	_, err := teamUseCase.userRepo.GetByID(ownerId)
	if err != nil {
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
	usr, err := teamUseCase.teamRepo.GetByID(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (teamUseCase *TeamUseCaseImpl) GetTeamsByNamePart(namePart string, limit uint) (models.Teams, error) {
	users, err := teamUseCase.teamRepo.GetTeamsByNamePart(namePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (teamUseCase *TeamUseCaseImpl) InviteMember(tid uint, uid uint, role usecases.Role) error {
	user, err := teamUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = teamUseCase.teamRepo.InviteMember(tid, user, role)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
