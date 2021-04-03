package usecases

import (
	"sport4all/app/models"
)

type SportUseCase interface {
	Create(sport *models.Sport) error
	GetAll(limit uint) (*[]models.Sport, error)
}
