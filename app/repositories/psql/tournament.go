package psql

import (
	"time"

	"github.com/jinzhu/gorm"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type TournamentStore struct {
	DB *gorm.DB
}

func CreateTournamentRepository(db *gorm.DB) repositories.TournamentRepository {
	return &TournamentStore{DB: db}
}

func (tournamentStore *TournamentStore) Create(tournament *models.Tournament) error {
	tournament.Created = time.Now().Unix()
	tournament.Status = models.NotStartedEvent

	if err := tournamentStore.DB.Create(tournament).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (tournamentStore *TournamentStore) GetByID(tid uint) (*models.Tournament, error) {
	tournament := new(models.Tournament)
	if err := tournamentStore.DB.Where("id = ?", tid).First(&tournament).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}

	return tournament, nil
}

func (tournamentStore *TournamentStore) GetTournamentByUser(uid uint) (*models.Tournaments, error) {
	ownerTournaments := new(models.Tournaments)
	if err := tournamentStore.DB.Model(&models.User{ID: uid}).
		Related(&ownerTournaments, "owner_id").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}

	return ownerTournaments, nil
}

func (tournamentStore *TournamentStore) Update(tournament *models.Tournament) error {
	oldTournament, err := tournamentStore.GetByID(tournament.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if tournament.Name != "" {
		oldTournament.Name = tournament.Name
	}
	if tournament.Location != "" {
		oldTournament.Location = tournament.Location
	}
	if tournament.Status != models.UnknownEvent {
		oldTournament.Status = tournament.Status
	}
	if tournament.System != "" {
		oldTournament.System = tournament.System
	}
	if tournament.About != "" {
		oldTournament.About = tournament.About
	}

	if err = tournamentStore.DB.Save(oldTournament).Error; err != nil {
		logger.Error(err)
		return errors.ErrInternal
	}

	return nil
}

func (tournamentStore *TournamentStore) AddTeam(tournamentId uint, teamId uint) error {
	if err := tournamentStore.DB.Model(&models.Tournament{ID: tournamentId}).
		Association("Teams").Append(&models.Team{ID: teamId}).Error; err != nil {
		logger.Error(err)
		return errors.ErrTournamentNotFound
	}

	return nil
}

func (tournamentStore *TournamentStore) RemoveTeam(tournamentId uint, teamId uint) error {
	if err := tournamentStore.DB.Model(&models.Tournament{ID: tournamentId}).
		Association("Teams").Delete(&models.Team{ID: teamId}).Error; err != nil {
		logger.Error(err)
		return errors.ErrTournamentNotFound
	}

	return nil
}

func (tournamentStore *TournamentStore) GetAllTeams(tournamentId uint) (*models.Teams, error) {
	var tournamentTeams models.Teams
	if err := tournamentStore.DB.Model(&models.Tournament{ID: tournamentId}).
		Related(&tournamentTeams, "Teams").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}
	return &tournamentTeams, nil
}

func (tournamentStore *TournamentStore) GetAllMeetings(tournamentId uint) (*models.Meetings, error) {
	var tournamentMeetings models.Meetings
	if err := tournamentStore.DB.Model(&models.Tournament{ID: tournamentId}).Where("next_meeting_id is null").
		Preload("PrevMeetings").Related(&tournamentMeetings, "tournamentId").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}

	return &tournamentMeetings, nil
}
