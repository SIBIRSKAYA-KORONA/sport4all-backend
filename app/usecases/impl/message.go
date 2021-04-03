package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/logger"
)

type MessageUseCaseImpl struct {
	messageRepo repositories.MessageRepository
}

func CreateMessageUseCase(messageRepo repositories.MessageRepository) usecases.MessageUseCase {
	return &MessageUseCaseImpl{messageRepo: messageRepo}
}

func (messageUseCase *MessageUseCaseImpl) Create(messages *[]models.Message) error {
	if err := messageUseCase.messageRepo.Create(messages); err != nil {
		logger.Info(err)
		return err
	}

	return nil
}

func (messageUseCase *MessageUseCaseImpl) GetAll(uid uint) (*[]models.Message, bool) {
	messages, has := messageUseCase.messageRepo.GetAll(uid)
	if !has {
		logger.Info("no messages for the user", uid)
		return nil, false
	}
	return messages, true
}

func (messageUseCase *MessageUseCaseImpl) UpdateAll(uid uint) error {
	if err := messageUseCase.messageRepo.UpdateAll(uid); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (messageUseCase *MessageUseCaseImpl) DeleteAll(uid uint) error {
	if err := messageUseCase.messageRepo.DeleteAll(uid); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
