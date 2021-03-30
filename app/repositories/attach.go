package repositories

import (
	"sport4all/app/models"
)

type AttachRepository interface {
	Create(attach *models.Attach) error
	Delete(key string) error
}