package usecases

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)


// TODO: move it for use cases
const (
	Unknown uint = iota
	New
	Progress
	Finished
)

type MeetingUseCase interface {
	Create(meeting *models.Meeting) error
	GetByID(mid uint) (*models.Meeting, error)
	Update(meeting *models.Meeting) error
}
