package usecases

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type MeetingStatUseCase interface {
	Create(stat *models.MeetingStat) error
	GetByID(statId uint) (*models.MeetingStat, error)
	Update(stat *models.MeetingStat) error
}
