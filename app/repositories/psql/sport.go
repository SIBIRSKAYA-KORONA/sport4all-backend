package psql

import (
	"github.com/jinzhu/gorm"
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type SportStore struct {
	db *gorm.DB
}

func CreateSportRepository(db *gorm.DB) repositories.SportRepository {
	return &SportStore{db: db}
}

func (sportStore *SportStore) Create(meeting *models.Sport) error {
	if err := sportStore.db.Create(meeting).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (sportStore *SportStore) Get(sportName string) (*models.Sport, error) {
	var sport models.Sport
	if err := sportStore.db.Where("name = ?", sportName).First(&sport).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return &sport, nil
}

func (sportStore *SportStore) GetAll(limit uint) (*[]models.Sport, error) {
	var sports []models.Sport
	if err := sportStore.db.Limit(limit).Find(&sports).Order("id").Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return &sports, nil
}
