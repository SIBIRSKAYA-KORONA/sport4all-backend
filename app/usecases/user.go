package usecases

import (
	"sport4all/app/models"
)

type UserUseCase interface {
	Create(user *models.User) (*models.Session, error)
	GetByID(uid uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	GetUserStats(uid uint) ([]models.Stats, error)
}
