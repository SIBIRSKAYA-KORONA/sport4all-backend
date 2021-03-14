package repositories

import (
	"sport4all/app/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(uid uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	IsValidPassword(password string, hashPassword []byte) bool
}
