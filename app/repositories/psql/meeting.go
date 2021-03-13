package psql

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"

	"github.com/jinzhu/gorm"
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

	if meeting.Status != models.Unknown {
		oldMeeting.Status = meeting.Status
	}
	if meeting.PrevGame != nil {
		oldMeeting.PrevGame = meeting.PrevGame
	}
	if meeting.NextGame != nil {
		oldMeeting.NextGame = meeting.NextGame
	}

	if err = meetingStore.DB.Save(oldMeeting).Error; err != nil {
		logger.Error(err)
		return errors.ErrInternal
	}

	return nil
}
