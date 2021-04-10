package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type MeetingUseCaseImpl struct {
	meetingRepo    repositories.MeetingRepository
	tournamentRepo repositories.TournamentRepository
}

func CreateMeetingUseCase(meetingRepo repositories.MeetingRepository,
	tournamentRepo repositories.TournamentRepository) usecases.MeetingUseCase {
	return &MeetingUseCaseImpl{meetingRepo: meetingRepo, tournamentRepo: tournamentRepo}
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
	old, err := meetingUseCase.GetByID(meeting.ID)
	if err != nil {
		return err
	}
	if old.Status > models.NotStartedEvent && meeting.Status <= old.Status {
		return errors.ErrMeetingStatusNotAcceptable
	}

	switch meeting.Status {
	case models.UnknownEvent, models.NotStartedEvent, models.RegistrationEvent:
		if err = meetingUseCase.meetingRepo.Update(meeting); err != nil {
			logger.Error(err)
			return err
		}
	case models.InProgressEvent, models.FinishedEvent:
		if err = meetingUseCase.meetingRepo.Update(&models.Meeting{ID: meeting.ID, Status: meeting.Status}); err != nil {
			logger.Error(err)
			return err
		}
		if meeting.Status == models.InProgressEvent {
			return nil
		}

		tournament, err := meetingUseCase.tournamentRepo.GetByID(old.TournamentId)
		if err != nil {
			logger.Warn(err)
			return err
		}
		if tournament.System == models.OlympicSystem && meeting.NextMeetingID != nil {
			stat, err := meetingUseCase.meetingRepo.GetMeetingTeamStat(meeting.ID)
			if err != nil || len(*stat) != 2 {
				logger.Warn(err)
				return err
			}
			winnerTeamId := (*stat)[0].TeamId
			if (*stat)[0].Score < (*stat)[1].Score {
				winnerTeamId = (*stat)[1].TeamId
			}
			if err = meetingUseCase.meetingRepo.AssignTeam(*meeting.NextMeetingID, winnerTeamId); err != nil {
				logger.Warn(err)
				return err
			}
		}
	default:
		return errors.ErrMeetingStatusNotAcceptable
	}

	return nil
}

func (meetingUseCase *MeetingUseCaseImpl) AssignTeam(mid uint, tid uint) error {
	meeting, err := meetingUseCase.meetingRepo.GetByID(mid)
	if err != nil {
		logger.Error(err)
		return err
	}

	result, err := meetingUseCase.tournamentRepo.IsTeamInTournament(meeting.TournamentId, tid)
	if err != nil {
		logger.Error(err)
		return err
	}

	if !result {
		return errors.ErrTeamNotFound
	}

	if err = meetingUseCase.meetingRepo.AssignTeam(mid, tid); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (meetingUseCase *MeetingUseCaseImpl) IsTeamInMeeting(mid uint, tid uint) (bool, error) {
	result, err := meetingUseCase.meetingRepo.IsTeamInMeeting(mid, tid)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	return result, nil
}

func (meetingUseCase *MeetingUseCaseImpl) UpdateTeamStat(stat *models.Stats) error {
	if err := meetingUseCase.meetingRepo.UpdateTeamStat(stat); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (meetingUseCase *MeetingUseCaseImpl) GetMeetingStat(mid uint) (*[]models.Stats, error) {
	stats, err := meetingUseCase.meetingRepo.GetMeetingStat(mid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return stats, nil
}
