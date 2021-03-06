package psql

import (
	"time"

	"github.com/jinzhu/gorm"
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type MeetingStore struct {
	db *gorm.DB
}

func CreateMeetingRepository(db *gorm.DB) repositories.MeetingRepository {
	return &MeetingStore{db: db}
}

func (meetingStore *MeetingStore) Create(meeting *models.Meeting) error {
	meeting.Status = models.RegistrationEvent
	if err := meetingStore.db.Create(meeting).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (meetingStore *MeetingStore) CreateBatch(meetings *[]models.Meeting) error {
	if err := meetingStore.db.Create(meetings).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (meetingStore *MeetingStore) GetByID(mid uint) (*models.Meeting, error) {
	meeting := new(models.Meeting)
	if err := meetingStore.db.Where("id = ?", mid).
		Preload("Attachments").
		Preload("Protocol").
		Preload("Teams.Players").
		Preload("Teams.Avatar").
		First(&meeting).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	return meeting, nil
}

func (meetingStore *MeetingStore) Update(meeting *models.Meeting) error {
	oldMeeting, err := meetingStore.GetByID(meeting.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if meeting.Status != models.UnknownEvent {
		oldMeeting.Status = meeting.Status
	}
	if meeting.Round != 0 {
		oldMeeting.Round = meeting.Round
	}
	if meeting.Group != 0 {
		oldMeeting.Group = meeting.Group
	}
	if meeting.TournamentId != 0 {
		oldMeeting.TournamentId = meeting.TournamentId
	}

	if err = meetingStore.db.Save(oldMeeting).Error; err != nil {
		logger.Error(err)
		return errors.ErrMeetingNotFound
	}

	return nil
}

func (meetingStore *MeetingStore) AssignTeam(mid uint, tid uint) error {
	var teams []models.Team

	if err := meetingStore.db.Model(&models.Meeting{ID: mid}).
		Association("teams").
		Find(&teams).
		Error; err != nil {
		logger.Error(err)
		return errors.ErrMeetingNotFound
	}

	if len(teams) >= 2 {
		logger.Error("Can't assign more than 2 teams on metting")
		return errors.ErrMeetingTeamsAreAlreadyAssigned
	}

	if err := meetingStore.db.Model(&models.Meeting{ID: mid}).
		Association("teams").
		Append(models.Team{ID: tid}).
		Error; err != nil {
		logger.Error(err)
		return errors.ErrTeamNotFound
	}

	return nil
}

func (meetingStore *MeetingStore) IsTeamInMeeting(mid uint, tid uint) (bool, error) {
	teams := new(models.Teams)
	if err := meetingStore.db.Model(&models.Meeting{ID: mid}).
		Related(&teams, "teams").Error; err != nil {
		logger.Error(err)
		return false, errors.ErrMeetingNotFound
	}

	for _, team := range *teams {
		if team.ID == tid {
			return true, nil
		}
	}

	return false, nil
}

func (meetingStore *MeetingStore) CreateStat(stat *models.Stats) error {
	var existsStat *models.Stats
	var err error

	if stat.PlayerId != nil {
		existsStat, err = meetingStore.GetMeetingPlayerStat(stat.MeetingId, stat.TeamId, *stat.PlayerId)
	} else {
		existsStat, err = meetingStore.GetMeetingTeamStat(stat.MeetingId, stat.TeamId)
	}

	if err != nil {
		stat.Created = time.Now().Unix()
		err = meetingStore.db.Create(stat).Error
	} else {
		stat.ID = existsStat.ID
		stat.Created = existsStat.Created
		err = meetingStore.db.Save(stat).Error
	}

	if err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (meetingStore *MeetingStore) CreatePlayersStats(stats *[]models.Stats) error {
	if err := meetingStore.db.Create(stats).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (meetingStore *MeetingStore) GetMeetingTeamsStats(mid uint) (*[]models.Stats, error) {
	var stats []models.Stats
	if err := meetingStore.db.Find(&stats, "player_id is null and meeting_id = ?", mid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	return &stats, nil
}

func (meetingStore *MeetingStore) GetMeetingPlayerStat(mid, tid, uid uint) (*models.Stats, error) {
	var stat models.Stats
	if err := meetingStore.db.First(&stat, "meeting_id = ? and team_id = ? and player_id = ?", mid, tid, uid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	return &stat, nil
}

func (meetingStore *MeetingStore) GetMeetingTeamStat(mid, tid uint) (*models.Stats, error) {
	var stat models.Stats
	if err := meetingStore.db.First(&stat, "meeting_id = ? and team_id = ? and player_id is null", mid, tid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	return &stat, nil
}

func (meetingStore *MeetingStore) GetMeetingStats(mid uint) (*[]models.Stats, error) {
	var stats []models.Stats
	if err := meetingStore.db.Model(&models.Meeting{ID: mid}).
		Related(&stats, "meetingId").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	return &stats, nil
}
