package usecases

import (
	"sport4all/app/models"
)

type SkillUseCase interface {
	Create(approvedUid, approvalUid uint, skill *models.Skill) error
	GetByNamePart(namePart string, limit uint) (*[]models.Skill, error)
	CreateApprove(approvedUid, approvalUid uint, approve *models.SkillApprove) error
	//Update(sid string) (uint, error)
	//Delete(sid string) error
}
