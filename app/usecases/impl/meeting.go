package impl

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type MeetingUseCaseImpl struct {
	meetingRepo repositories.MeetingRepository
}

func CreateMeetingUseCase(meetingRepo repositories.MeetingRepository) usecases.MeetingUseCase {
	return &MeetingUseCaseImpl{meetingRepo: meetingRepo}
}

func (meetingUseCase *MeetingUseCaseImpl) Create(meeting *models.Meeting) error {
	if err := meetingUseCase.meetingRepo.Create(meeting); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (meetingUseCase *MeetingUseCaseImpl) GetByID(mid uint) (*models.Meeting, error) {
	meeting, err := meetingUseCase.meetingRepo.GetByID(mid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return meeting, nil
}

func (meetingUseCase *MeetingUseCaseImpl) Update(meeting *models.Meeting) error {
	return nil
}
