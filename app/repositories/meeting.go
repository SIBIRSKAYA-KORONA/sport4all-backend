package repositories

import (
	"sport4all/app/models"
)

type MeetingRepository interface {
	Create(meeting *models.Meeting) error
	CreateBatch(meetings *[]models.Meeting) error
	GetByID(mid uint) (*models.Meeting, error)
	Update(meeting *models.Meeting) error
	AssignTeam(mid uint, tid uint) error
	IsTeamInMeeting(mid uint, tid uint) (bool, error)
	UpdateTeamStat(stat *models.Stats) error
	GetMeetingTeamStat(mid uint) (*[]models.Stats, error)
	GetMeetingStat(mid uint) (*[]models.Stats, error)
}
