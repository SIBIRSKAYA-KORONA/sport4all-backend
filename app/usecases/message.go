package usecases

import (
	"sport4all/app/models"
)

type MessageUseCase interface {
	Create(messages *[]models.Message) error
	GetAll(uid uint) (*[]models.Message, bool)
	DeleteAll(uid uint) error
	UpdateAll(uid uint) error
}
