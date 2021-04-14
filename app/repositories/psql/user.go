package psql

import (
	"strings"
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

func (userStore *UserStore) Update(user *models.User) error {
	var oldUser models.User
	if err := userStore.db.Where("id = ?", user.ID).First(&oldUser).Error; err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	if user.Name != "" {
		oldUser.Name = user.Name
	}
	if user.Surname != "" {
		oldUser.Surname = user.Surname
	}
	if user.Nickname != "" {
		oldUser.Nickname = user.Nickname
	}
	if user.About != "" {
		oldUser.About = user.About
	}
	if err := userStore.db.Save(oldUser).Error; err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	return nil
}

func (userStore *UserStore) GetByID(uid uint) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.db.Where("id = ?", uid).Preload("Avatar").First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return usr, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.db.Where("nickname = ?", nickname).Preload("Avatar").First(&usr).Error; err != nil {
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
		Preload("Approvals", "to_uid = ?", uid).
		Related(&skills, "skills").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	for i := range skills {
		for j := range skills[i].Approvals {
			skills[i].Approvals[j].FromUser, _ = userStore.GetByID(skills[i].Approvals[j].FromUid)
		}
	}

	return &skills, nil
}

func (userStore *UserStore) GetUserStats(uid uint) (*[]models.Stats, error) {
	var stats []models.Stats
	if err := userStore.db.Model(&models.User{ID: uid}).Order("id desc").
		Related(&stats, "playerId").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return &stats, nil
}

func (userStore *UserStore) SearchUsers(uid *uint, namePart string, limit uint) (*[]models.User, error) {
	var users []models.User

	lowerName := "%" + strings.ToLower(namePart) + "%"
	query := userStore.db.Select("id, name, surname, nickname").
		Order("name, surname, nickname").Limit(limit).
		Where("LOWER(name) LIKE ?", lowerName).
		Or("LOWER(surname) LIKE ?", lowerName).
		Or("LOWER(nickname) LIKE ?", lowerName).
		Preload("Avatar")

	if uid != nil {
		query = query.Not("id", *uid)
	}

	if err := query.Find(&users).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return &users, nil
}
