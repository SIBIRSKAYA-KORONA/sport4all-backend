package usecases

import (
	"sport4all/app/models"
)

type MeetingUseCase interface {
	Create(meeting *models.Meeting) error
	GetByID(mid uint) (*models.Meeting, error)
	Update(meeting *models.Meeting) error
}
