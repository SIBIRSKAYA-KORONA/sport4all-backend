package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type SessionUseCaseImpl struct {
	sessionRepo repositories.SessionRepository
	userRepo    repositories.UserRepository
}

func CreateSessionUseCase(sessionRepo repositories.SessionRepository, userRepo repositories.UserRepository) usecases.SessionUseCase {
	return &SessionUseCaseImpl{sessionRepo: sessionRepo, userRepo: userRepo}
}

func (sessionUseCaseImpl *SessionUseCaseImpl) Create(user *models.User) (*models.Session, error) {
	realUser, err := sessionUseCaseImpl.userRepo.GetByNickname(user.Nickname)
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	if !sessionUseCaseImpl.userRepo.IsValidPassword(user.Password, realUser.HashPassword) {
		logger.Error(errors.ErrWrongPassword)
		return nil, errors.ErrWrongPassword
	}

	ses := &models.Session{ID: realUser.ID}
	if err = sessionUseCaseImpl.sessionRepo.Create(ses); err != nil {
		logger.Error(err)
		return nil, err
	}
	return ses, nil
}

func (sessionUseCaseImpl *SessionUseCaseImpl) GetByID(sid string) (uint, error) {
	uid, err := sessionUseCaseImpl.sessionRepo.Get(sid)
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	return uid, nil
}

func (sessionUseCaseImpl *SessionUseCaseImpl) Delete(sid string) error {
	if err := sessionUseCaseImpl.sessionRepo.Delete(sid); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
