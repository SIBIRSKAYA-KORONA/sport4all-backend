package usecases

import (
	"sport4all/app/models"
)

type UserUseCase interface {
	Create(user *models.User) (*models.Session, error)
	Update(user *models.User) error
	GetByID(uid uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	GetUserSkills(uid uint) (*[]models.Skill, error)
	GetUserStats(uid uint) (*[]models.Stats, error)
	SearchUsers(sid string, namePart string, limit uint) (*[]models.User, error)
}
