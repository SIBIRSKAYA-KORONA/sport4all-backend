package psql

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type SkillStore struct {
	db *gorm.DB
}

func CreateSkillRepository(db *gorm.DB) repositories.SkillRepository {
	return &SkillStore{db: db}
}

func (skillStore *SkillStore) Create(toUid, fromUid uint, skill *models.Skill) error {
	if err := skillStore.db.Create(skill).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	approve := models.SkillApprove{
		SkillId: skill.ID,
		FromUid: fromUid,
		ToUid:   toUid,
	}
	if err := skillStore.CreateApprove(&approve); err != nil {
		logger.Warn(err)
	}
	skill.Approvals = append(skill.Approvals, approve)

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

func (skillStore *SkillStore) CreateApprove(approve *models.SkillApprove) error {
	if err := skillStore.db.Create(approve).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	approve.CreateAt = time.Now().Unix()
	if err := skillStore.db.Model(&models.Skill{ID: approve.SkillId}).
		Association("users").
		Append(models.User{ID: approve.ToUid}).
		Error; err != nil {
		logger.Warn(err)
	}

	return nil
}
