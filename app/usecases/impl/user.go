package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/logger"
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

func (userUseCase *UserUseCaseImpl) Update(user *models.User) error {
	err := userUseCase.userRepo.Update(user)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
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

func (userUseCase *UserUseCaseImpl) GetUserSkills(uid uint) (*[]models.Skill, error) {
	skill, err := userUseCase.userRepo.GetUserSkills(uid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return skill, nil
}

func (userUseCase *UserUseCaseImpl) GetUserStats(uid uint) (*[]models.Stats, error) {
	stats, err := userUseCase.userRepo.GetUserStats(uid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return stats, nil
}

func (userUseCase *UserUseCaseImpl) SearchUsers(uid *uint, nicknamePart string, limit uint) (*[]models.User, error) {
	users, err := userUseCase.userRepo.SearchUsers(uid, nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return users, nil
}
