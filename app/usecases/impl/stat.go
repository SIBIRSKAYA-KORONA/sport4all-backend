package impl

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type MeetingStatUseCaseImpl struct {
	meetingRepo repositories.MeetingStatRepository
}

func CreateMeetingStatUseCase(meetingRepo repositories.MeetingStatRepository) usecases.MeetingStatUseCase {
	return &MeetingStatUseCaseImpl{meetingRepo: meetingRepo}
}

func (meetingStatUseCase *MeetingStatUseCaseImpl) Create(meeting *models.MeetingStat) error {
	if err := meetingStatUseCase.meetingRepo.Create(meeting); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (meetingStatUseCase *MeetingStatUseCaseImpl) GetByID(mid uint) (*models.MeetingStat, error) {
	meeting, err := meetingStatUseCase.meetingRepo.GetByID(mid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return meeting, nil
}

func (meetingStatUseCase *MeetingStatUseCaseImpl) Update(meeting *models.MeetingStat) error {
	return nil
}
