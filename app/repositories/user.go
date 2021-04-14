package repositories

import (
	"sport4all/app/models"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	GetByID(uid uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	IsValidPassword(password string, hashPassword []byte) bool
	GetUserSkills(uid uint) (*[]models.Skill, error)
	GetUserStats(uid uint) (*[]models.Stats, error)
	SearchUsers(uid *uint, namePart string, limit uint) (*[]models.User, error)
}
