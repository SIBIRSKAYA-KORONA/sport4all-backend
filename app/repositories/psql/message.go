package psql

import (
	"github.com/jinzhu/gorm"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type MessageStore struct {
	db *gorm.DB
}

func CreateMessageRepository(db *gorm.DB) repositories.MessageRepository {
	return &MessageStore{db: db}
}

func (messageStore *MessageStore) Create(messages *[]models.Message) error {
	for _, message := range *messages {
		if err := messageStore.db.Create(&message).Error; err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func (messageStore *MessageStore) UpdateAll(uid uint) error {
	err := messageStore.db.Model(models.Message{}).
		Where("target_uid = ?", uid).UpdateColumn("is_read", true).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	return nil
}

func (messageStore *MessageStore) DeleteAll(uid uint) error {
	err := messageStore.db.Where("target_uid = ?", uid).Delete(models.Message{}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	return nil
}

func (messageStore *MessageStore) GetAll(uid uint) (*[]models.Message, bool) {
	var messages []models.Message
	err := messageStore.db.Order("id desc").Where("target_uid = ?", uid).Find(&messages).Error
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return &messages, true
}
