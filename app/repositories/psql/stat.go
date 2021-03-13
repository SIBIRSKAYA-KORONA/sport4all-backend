package psql

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"

	"github.com/jinzhu/gorm"
)

type MeetingStatStore struct {
	DB *gorm.DB
}

func CreateMeetingStatRepository(db *gorm.DB) repositories.MeetingStatRepository {
	return &MeetingStatStore{DB: db}
}

func (meetingStatStore *MeetingStatStore) Create(meeting *models.MeetingStat) error {
	if err := meetingStatStore.DB.Create(meeting).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (meetingStatStore *MeetingStatStore) GetByID(mid uint) (*models.MeetingStat, error) {
	meeting := new(models.MeetingStat)
	if err := meetingStatStore.DB.Where("id = ?", mid).First(&meeting).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrMeetingNotFound
	}

	return meeting, nil
}

func (meetingStatStore *MeetingStatStore) Update(meeting *models.MeetingStat) error {
	_, err := meetingStatStore.GetByID(meeting.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	//if meeting.Status != models.Unknown {
	//	oldMeeting.Status = meeting.Status
	//}
	//if meeting.PrevGame != nil {
	//	oldMeeting.PrevGame = meeting.PrevGame
	//}
	//if meeting.NextGame != nil {
	//	oldMeeting.NextGame = meeting.NextGame
	//}
	//
	//if err = meetingStore.DB.Save(oldMeeting).Error; err != nil {
	//	logger.Error(err)
	//	return errors.ErrInternal
	//}

	return nil
}
