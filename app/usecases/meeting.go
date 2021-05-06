package usecases

import (
	"sport4all/app/models"
)

type MeetingUseCase interface {
	Create(meeting *models.Meeting) error
	GetByID(mid uint) (*models.Meeting, error)
	Update(meeting *models.Meeting) error
	AssignTeam(mid uint, tid uint) error
	IsTeamInMeeting(mid uint, tid uint) (bool, error)
	UpdateTeamStat(stat *models.Stats) error
	GetMeetingStat(mid uint) (*[]models.Stats, error)
	GetStatsByImage(mid uint, imagePath, protocolType string) (*[]models.Stats, error)
}
