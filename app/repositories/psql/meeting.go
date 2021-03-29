package psql

import (
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
	if err := meetingStore.db.Where("id = ?", mid).First(&meeting).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	if err := meetingStore.db.Where("meeting_id = ?", mid).Find(&meeting.Attachments).Error; err != nil {
		logger.Warn("meeting attachments not found: ", err)
	}

	if err := meetingStore.db.Model(meeting).
		Related(&meeting.Teams, "teams").
		Order("id").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrInternal
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
		return errors.ErrInternal
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
		return errors.ErrInternal
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

func (meetingStore *MeetingStore) UpdateTeamStat(stat *models.Stats) error {
	if err := meetingStore.db.Create(stat).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (meetingStore *MeetingStore) GetMeetingStat(mid uint) (*[]models.Stats, error) {
	var stats []models.Stats
	if err := meetingStore.db.Model(&models.Meeting{ID: mid}).
		Related(&stats, "meetingId").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	return &stats, nil
}
