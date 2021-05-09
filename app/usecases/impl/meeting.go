package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/math"
	"strconv"
	"strings"
)

type MeetingUseCaseImpl struct {
	meetingRepo    repositories.MeetingRepository
	tournamentRepo repositories.TournamentRepository
	ocrRepo        repositories.OcrRepository
}

func CreateMeetingUseCase(meetingRepo repositories.MeetingRepository,
	tournamentRepo repositories.TournamentRepository, ocrRepo repositories.OcrRepository) usecases.MeetingUseCase {
	return &MeetingUseCaseImpl{meetingRepo: meetingRepo, tournamentRepo: tournamentRepo, ocrRepo: ocrRepo}
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
	case models.InProgressEvent:
		if err = meetingUseCase.meetingRepo.Update(&models.Meeting{ID: meeting.ID, Status: meeting.Status}); err != nil {
			logger.Error(err)
			return err
		}
	case models.FinishedEvent:
		stat, err := meetingUseCase.meetingRepo.GetMeetingTeamStat(meeting.ID)
		if err != nil {
			logger.Error(err)
			return errors.ErrMeetingStatusNotAcceptable
		}
		if err = meetingUseCase.meetingRepo.Update(&models.Meeting{ID: meeting.ID, Status: meeting.Status}); err != nil {
			logger.Error(err)
			return err
		}
		tournament, err := meetingUseCase.tournamentRepo.GetByID(old.TournamentId)
		if err != nil {
			logger.Warn(err)
			return err
		}
		if tournament.System == models.OlympicSystem && old.NextMeetingID != nil {
			winnerTeamId := (*stat)[0].TeamId
			if (*stat)[0].Score < (*stat)[1].Score {
				winnerTeamId = (*stat)[1].TeamId
			}
			if err = meetingUseCase.meetingRepo.AssignTeam(*old.NextMeetingID, winnerTeamId); err != nil {
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

func (meetingUseCase *MeetingUseCaseImpl) CreateTeamStat(stat *models.Stats) error {
	if err := meetingUseCase.meetingRepo.CreateTeamStat(stat); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (meetingUseCase *MeetingUseCaseImpl) calcTeamStats(stats *[]models.Stats) error {
	if len(*stats) == 0 {
		return nil
	}

	teamStats := make(map[uint]uint, 0)
	for idx := range *stats {
		stat := (*stats)[idx]
		if stat.PlayerId == nil {
			return errors.ErrInvalidStats
		}
		teamStats[stat.TeamId] += stat.Score
	}
	mid := (*stats)[0].MeetingId

	for teamId, score := range teamStats {
		stat := models.Stats{
			Score:     score,
			MeetingId: mid,
			TeamId:    teamId,
			PlayerId:  nil,
		}
		*stats = append(*stats, stat)
	}

	return nil
}

func (meetingUseCase *MeetingUseCaseImpl) CreatePlayersStats(stats *[]models.Stats) error {
	err := meetingUseCase.calcTeamStats(stats)
	if err != nil {
		logger.Error(err)
		return err
	}

	for idx := range *stats {
		stat := (*stats)[idx]
		if err = meetingUseCase.CreateTeamStat(&stat); err != nil {
			logger.Error(err)
			return err
		}
	}

	/*
		TODO: don't work batch create
		if err := meetingUseCase.meetingRepo.CreatePlayersStats(stats); err != nil {
			logger.Error(err)
			return err
		}
	*/

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

func (meetingUseCase *MeetingUseCaseImpl) GetStatsByImage(mid uint,
	imagePath, protocolType string) (*[]models.Stats, error) {

	protocolImage := models.ProtocolImage{Path: imagePath}

	info, has := models.ProtocolTypes[protocolType]
	if !has {
		// may be unknown format (1,6)
		arr := strings.Split(protocolType, ",")
		if len(arr) < 2 {
			return nil, errors.ErrProtocolTypeNotFound
		}

		tmp, err := strconv.ParseInt(arr[0], 10, 32)
		if err != nil {
			return nil, errors.ErrProtocolTypeNotFound
		}
		info.PlayerColumn = int32(tmp)

		tmp, err = strconv.ParseInt(arr[1], 10, 32)
		if err != nil {
			return nil, errors.ErrProtocolTypeNotFound
		}
		info.ScoreColumn = int32(tmp)
	}
	protocolImage.Info = info

	meeting, err := meetingUseCase.GetByID(mid)
	if err != nil {
		return nil, err
	}

	extractedStats, err := meetingUseCase.ocrRepo.GetStatsByImage(&protocolImage)
	if err != nil {
		return nil, err
	}

	stats := make([]models.Stats, 0)
	for _, extractedStat := range *extractedStats {
		for _, team := range meeting.Teams {
			for idx := range team.Players {
				player := team.Players[idx]
				if math.LevenshteinDist(player.Name, extractedStat.Name) < 2 &&
					math.LevenshteinDist(player.Surname, extractedStat.Surname) < 2 {
					stats = append(stats,
						models.Stats{
							Score:     uint(extractedStat.Score),
							MeetingId: mid,
							TeamId:    team.ID,
							PlayerId:  &player.ID,
						})
				}
			}
		}
	}

	if err = meetingUseCase.calcTeamStats(&stats); err != nil {
		logger.Error(err)
		return nil, err
	}

	return &stats, nil
}
