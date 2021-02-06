package impl

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type UserUseCaseImpl struct {
	sessionRepo repositories.SessionRepository
	userRepo    repositories.UserRepository
}

func CreateUserUseCase(sessionRepo repositories.SessionRepository, userRepo repositories.UserRepository) usecases.UserUseCase {
	return &UserUseCaseImpl{sessionRepo: sessionRepo, userRepo: userRepo}
}

func (userUseCase *UserUseCaseImpl) Create(user *models.User) (*models.Session, error) {
	if err := userUseCase.userRepo.Create(user); err != nil {
		logger.Error(err)
		return nil, err
	}

	ses := &models.Session{ID: user.ID}
	if err := userUseCase.sessionRepo.Create(ses); err != nil {
		logger.Error(err)
		return nil, err
	}
	return ses, nil
}

func (userUseCase *UserUseCaseImpl) GetByID(uid uint) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (userUseCase *UserUseCaseImpl) GetByNickname(nickname string) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByNickname(nickname)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}
