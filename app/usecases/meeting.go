package usecases

import (
	"sport4all/app/models"
)

// TODO: move it for all use cases
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
