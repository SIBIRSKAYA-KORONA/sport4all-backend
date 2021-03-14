package repositories

import (
	"sport4all/app/models"
)

type MeetingRepository interface {
	Create(meeting *models.Meeting) error
	GetByID(mid uint) (*models.Meeting, error)
	Update(meeting *models.Meeting) error
}
