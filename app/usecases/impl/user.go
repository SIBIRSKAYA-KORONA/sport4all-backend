package impl

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type UserUseCaseImpl struct {
	userRepo repositories.UserRepository
}

func CreateUserUseCase(userRepo_ repositories.UserRepository) usecases.UserUseCase {
	return &UserUseCaseImpl{
		userRepo: userRepo_,
	}
}

func (userUseCase *UserUseCaseImpl) Create(user *models.User) error {
	err := userUseCase.userRepo.Create(user)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (userUseCase *UserUseCaseImpl) GetByNickname(nickname string) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByNickname(nickname)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}
