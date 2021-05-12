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
	CreateStat(stat *models.Stats) error
	CreatePlayersStats(stats *[]models.Stats) error
	GetMeetingTeamsStats(mid uint) (*[]models.Stats, error)
	GetMeetingPlayerStat(mid, tid, uid uint) (*models.Stats, error)
	GetMeetingStats(mid uint) (*[]models.Stats, error)
}
