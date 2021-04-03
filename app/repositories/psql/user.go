package psql

import (
	"time"

	"github.com/jinzhu/gorm"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/hasher"
	"sport4all/pkg/logger"
)

type UserStore struct {
	db *gorm.DB
}

func CreateUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserStore{db: db}
}

func (userStore *UserStore) Create(usr *models.User) error {
	usr.Created = time.Now().Unix()
	usr.HashPassword = hasher.HashPassword(usr.Password)
	if err := userStore.db.Create(usr).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (userStore *UserStore) GetByID(uid uint) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.db.Where("id = ?", uid).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	if err := userStore.db.Where("user_id = ?", uid).First(&usr.Avatar).Error; err != nil {
		logger.Warn("user avatar not found: ", err)
	}

	return usr, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.db.Where("nickname = ?", nickname).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return usr, nil
}

func (userStore *UserStore) IsValidPassword(password string, hashPassword []byte) bool {
	return hasher.IsEqualPassword(password, hashPassword)
}

func (userStore *UserStore) GetUserSkills(uid uint) (*[]models.Skill, error) {
	var skills []models.Skill
	if err := userStore.db.Model(&models.User{ID: uid}).
		Preload("Approvals").
		Preload("Approvals.Users").
		Related(&skills, "Skills").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return &skills, nil
}

func (userStore *UserStore) GetUserStats(uid uint) (*[]models.Stats, error) {
	var stats []models.Stats
	if err := userStore.db.Model(&models.User{ID: uid}).
		Related(&stats, "playerId").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return &stats, nil
}
