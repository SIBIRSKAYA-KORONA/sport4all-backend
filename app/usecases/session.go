package usecases

import (
	"sport4all/app/models"
)

type SessionUseCase interface {
	Create(user *models.User) (*models.Session, error)
	GetByID(sid string) (uint, error)
	Delete(sid string) error
}
