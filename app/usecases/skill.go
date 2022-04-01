package usecases

import (
	"sport4all/app/models"
)

type SkillUseCase interface {
	Create(toUid, fromUid uint, skill *models.Skill) error
	GetByNamePart(namePart string, limit uint) (*[]models.Skill, error)
	CreateApprove(approve *models.SkillApprove) error
	// Update(sid string) (uint, error)
	// Delete(sid string) error
}
