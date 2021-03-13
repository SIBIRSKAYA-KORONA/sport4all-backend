package repositories

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type MeetingStatRepository interface {
	Create(meeting *models.MeetingStat) error
	GetByID(mid uint) (*models.MeetingStat, error)
	Update(meeting *models.MeetingStat) error
}
