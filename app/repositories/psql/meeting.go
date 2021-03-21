package psql

import (
	"github.com/jinzhu/gorm"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type MeetingStore struct {
	DB *gorm.DB
}

func CreateMeetingRepository(db *gorm.DB) repositories.MeetingRepository {
	return &MeetingStore{DB: db}
}

func (meetingStore *MeetingStore) Create(meeting *models.Meeting) error {
	if err := meetingStore.DB.Create(meeting).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (meetingStore *MeetingStore) GetByID(mid uint) (*models.Meeting, error) {
	meeting := new(models.Meeting)
	if err := meetingStore.DB.Where("id = ?", mid).First(&meeting).Error; err != nil {
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
	if meeting.Stats != "" {
		oldMeeting.Stats = meeting.Stats
	}
	if meeting.PrevMeetings != nil && len(meeting.PrevMeetings) != 0 {
		oldMeeting.PrevMeetings = meeting.PrevMeetings
	}
	if meeting.NextMeetingID != nil && *meeting.NextMeetingID != 0 {
		oldMeeting.NextMeetingID = meeting.NextMeetingID
	}

	if err = meetingStore.DB.Save(oldMeeting).Error; err != nil {
		logger.Error(err)
		return errors.ErrInternal
	}

	return nil
}
