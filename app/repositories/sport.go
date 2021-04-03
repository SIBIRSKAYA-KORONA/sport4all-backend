package repositories

import (
	"sport4all/app/models"
)

type SportRepository interface {
	Create(sport *models.Sport) error
	Get(sportName string) (*models.Sport, error)
	GetAll(limit uint) (*[]models.Sport, error)
}
