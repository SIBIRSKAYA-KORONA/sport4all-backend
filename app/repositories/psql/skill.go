package psql

import (
	"github.com/jinzhu/gorm"
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"strings"
)

type SkillStore struct {
	db *gorm.DB
}

func CreateSkillRepository(db *gorm.DB) repositories.SkillRepository {
	return &SkillStore{db: db}
}

func (skillStore *SkillStore) Create(skill *models.Skill) error {
	if err := skillStore.db.Create(skill).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (skillStore *SkillStore) GetByNamePart(namePart string, limit uint) (*[]models.Skill, error) {
	skills := make([]models.Skill, 0)
	if err := skillStore.db.Limit(limit).Where("LOWER(name) LIKE ?", strings.ToLower(namePart)+"%").
		Find(&skills).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrSkillNotFound
	}

	return &skills, nil
}
