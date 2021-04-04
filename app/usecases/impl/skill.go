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

func (skillUseCaseImpl *SkillUseCaseImpl) Create(toUid, fromUid uint, skill *models.Skill) error {
	if _, err := skillUseCaseImpl.userRepo.GetByID(toUid); err != nil {
		return err
	}

	if err := skillUseCaseImpl.skillRepo.Create(toUid, fromUid, skill); err != nil {
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

func (skillUseCaseImpl *SkillUseCaseImpl) CreateApprove(approve *models.SkillApprove) error {
	if _, err := skillUseCaseImpl.userRepo.GetByID(approve.ToUid); err != nil {
		return err
	}

	if err := skillUseCaseImpl.skillRepo.CreateApprove(approve); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
