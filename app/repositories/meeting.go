package repositories

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type MeetingRepository interface {
	Create(meeting *models.Meeting) error
	GetByID(mid uint) (*models.Meeting, error)
	Update(meeting *models.Meeting) error
}
