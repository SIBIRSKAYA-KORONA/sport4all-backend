package usecases

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type SessionUseCase interface {
	Create(user *models.User) (*models.Session, error)
	GetByID(sid string) (uint, error)
	Delete(sid string) error
}
