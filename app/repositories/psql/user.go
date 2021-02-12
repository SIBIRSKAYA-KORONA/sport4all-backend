package psql

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/hasher"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type UserStore struct {
	DB *gorm.DB
}

func CreateUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserStore{DB: db}
}

func (userStore *UserStore) Create(usr *models.User) error {
	usr.Created = time.Now().Unix()
	usr.HashPassword = hasher.HashPassword(usr.Password)
	if err := userStore.DB.Create(usr).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (userStore *UserStore) GetByID(uid uint) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.DB.Where("id = ?", uid).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}
	return usr, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.DB.Where("nickname = ?", nickname).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}
	return usr, nil
}

func (userStore *UserStore) IsValidPassword(password string, hashPassword []byte) bool {
	return hasher.IsEqualPassword(password, hashPassword)
}
