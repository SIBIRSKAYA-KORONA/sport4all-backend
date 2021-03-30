package usecases

import (
	"sport4all/app/models"
)

type AttachUseCase interface {
	Create(attach *models.Attach) error
	Delete(key string) error
}
