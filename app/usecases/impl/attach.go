package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/logger"
)

type AttachUseCaseImpl struct {
	attachRepo repositories.AttachRepository
}

func CreateAttachUseCase(attachRepo repositories.AttachRepository) usecases.AttachUseCase {
	return &AttachUseCaseImpl{attachRepo: attachRepo}
}

func (attachUseCase *AttachUseCaseImpl) Create(attach *models.Attach) error {
	if attach.MeetingId == nil {
		attachments, err := attachUseCase.attachRepo.GetByEntityID(attach)
		if err == nil {
			for _, value := range *attachments {
				_ = attachUseCase.attachRepo.Delete(value.Key)
			}
		}
	}

	if err := attachUseCase.attachRepo.Create(attach); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (attachUseCase *AttachUseCaseImpl) Delete(key string) error {
	if err := attachUseCase.attachRepo.Delete(key); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
