package repositories

import "sport4all/app/models"

type MessageRepository interface {
	Create(messages *[]models.Message) error
	GetAll(uid uint) (*[]models.Message, bool)
	DeleteAll(uid uint) error
	UpdateAll(uid uint) error
}
