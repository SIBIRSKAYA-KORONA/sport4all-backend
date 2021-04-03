package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/logger"
)

type SkillUseCaseImpl struct {
	skillRepo repositories.SkillRepository
	userRepo  repositories.UserRepository
}

func CreateSkillUseCase(skillRepo repositories.SkillRepository) usecases.SkillUseCase {
	return &SkillUseCaseImpl{skillRepo: skillRepo}
}

func (skillUseCaseImpl *SkillUseCaseImpl) Create(approvedUid, approvalUid uint, skill *models.Skill) error {
	if _, err := skillUseCaseImpl.userRepo.GetByID(approvedUid); err != nil {
		return err
	}

	skill.Users = append(skill.Users, models.User{ID: approvedUid})
	skill.Approvals = append(skill.Approvals, models.SkillApprove{SkillId: &skill.ID})
	skill.Approvals[0].Users = append(skill.Approvals[0].Users, models.User{ID: approvalUid})
	if err := skillUseCaseImpl.skillRepo.Create(skill); err != nil {
		logger.Info(err)
		return err
	}

	return nil
}

func (skillUseCaseImpl *SkillUseCaseImpl) GetByNamePart(namePart string, limit uint) (*[]models.Skill, error) {
	skills, err := skillUseCaseImpl.skillRepo.GetByNamePart(namePart, limit)
	if err != nil {
		logger.Info(err)
		return nil, err
	}

	return skills, nil
}
