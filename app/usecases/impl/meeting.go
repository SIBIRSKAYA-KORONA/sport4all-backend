package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/logger"
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
	// TODO: валидация состояния турнира
	if err := meetingUseCase.meetingRepo.Update(meeting); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
