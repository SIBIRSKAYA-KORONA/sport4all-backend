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

func CreateSkillUseCase(skillRepo repositories.SkillRepository, userRepo repositories.UserRepository) usecases.SkillUseCase {
	return &SkillUseCaseImpl{skillRepo: skillRepo, userRepo: userRepo}
}

func (skillUseCaseImpl *SkillUseCaseImpl) Create(approvedUid, approvalUid uint, skill *models.Skill) error {
	if _, err := skillUseCaseImpl.userRepo.GetByID(approvedUid); err != nil {
		return err
	}

	if err := skillUseCaseImpl.skillRepo.Create(approvedUid, approvalUid, skill); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (skillUseCaseImpl *SkillUseCaseImpl) GetByNamePart(namePart string, limit uint) (*[]models.Skill, error) {
	skills, err := skillUseCaseImpl.skillRepo.GetByNamePart(namePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return skills, nil
}

func (skillUseCaseImpl *SkillUseCaseImpl) CreateApprove(approvedUid, approvalUid uint, approve *models.SkillApprove) error {
	if _, err := skillUseCaseImpl.userRepo.GetByID(approvedUid); err != nil {
		return err
	}

	if err := skillUseCaseImpl.skillRepo.CreateApprove(approvedUid, approvalUid, approve); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
