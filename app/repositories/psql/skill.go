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

func (skillStore *SkillStore) Create(approvedUid, approvalUid uint, skill *models.Skill) error {
	skill.Approvals = append(skill.Approvals, models.SkillApprove{SkillId: &skill.ID})

	if err := skillStore.db.Create(skill).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	if err := skillStore.db.Model(&models.Skill{ID: skill.ID}).
		Association("users").
		Append(models.User{ID: approvedUid}).
		Error; err != nil {
		logger.Error(err)
		return errors.ErrSkillNotFound
	}

	if err := skillStore.db.Model(&models.SkillApprove{ID: skill.Approvals[0].ID}).
		Association("users").
		Append(models.User{ID: approvalUid}).
		Error; err != nil {
		logger.Error(err)
		return errors.ErrSkillNotFound
	}

	return nil
}

func (skillStore *SkillStore) GetByNamePart(namePart string, limit uint) (*[]models.Skill, error) {
	skills := make([]models.Skill, 0)
	if err := skillStore.db.Limit(limit).Where("LOWER(name) LIKE ?", "%"+strings.ToLower(namePart)+"%").
		Find(&skills).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrSkillNotFound
	}

	return &skills, nil
}

func (skillStore *SkillStore) CreateApprove(approvedUid, approvalUid uint, approve *models.SkillApprove) error {
	if err := skillStore.db.Create(approve).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	if err := skillStore.db.Model(&models.Skill{ID: *approve.SkillId}).
		Association("users").
		Append(models.User{ID: approvedUid}).
		Error; err != nil {
		logger.Warn(err)
	}

	if err := skillStore.db.Model(&models.SkillApprove{ID: approve.ID}).
		Association("users").
		Append(models.User{ID: approvalUid}).
		Error; err != nil {
		logger.Error(err)
		return errors.ErrSkillNotFound
	}

	return nil
}
