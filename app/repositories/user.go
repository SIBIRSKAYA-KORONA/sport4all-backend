package repositories

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByNickname(nickname string) (*models.User, error)
}
