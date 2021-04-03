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
	db *gorm.DB
}

func CreateTournamentRepository(db *gorm.DB) repositories.TournamentRepository {
	return &TournamentStore{db: db}
}

func (tournamentStore *TournamentStore) Create(tournament *models.Tournament) error {
	tournament.Created = time.Now().Unix()
	tournament.Status = models.NotStartedEvent

	if err := tournamentStore.db.Create(tournament).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (tournamentStore *TournamentStore) GetByID(tid uint) (*models.Tournament, error) {
	tournament := new(models.Tournament)
	if err := tournamentStore.db.Where("id = ?", tid).First(&tournament).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}

	if err := tournamentStore.db.Where("tournament_id = ?", tid).First(&tournament.Avatar).Error; err != nil {
		logger.Warn("tournament avatar not found: ", err)
	}

	return tournament, nil
}

func (tournamentStore *TournamentStore) GetTournamentByUserOwner(uid uint) (*models.Tournaments, error) {
	var ownerTournaments models.Tournaments
	if err := tournamentStore.db.Model(&models.User{ID: uid}).
		Related(&ownerTournaments, "owner_id").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}

	for idx := range ownerTournaments {
		if err := tournamentStore.db.Where("tournament_id = ?", ownerTournaments[idx].ID).
			First(&ownerTournaments[idx].Avatar).Error; err != nil {
			logger.Warn("tournament", ownerTournaments[idx].ID, " avatar not found: ", err)
		}
	}

	return &ownerTournaments, nil
}

func (tournamentStore *TournamentStore) IsTournamentOrganizer(tournamentID uint, userID uint) (bool, error) {
	tournament := new(models.Tournament)
	if err := tournamentStore.db.Where("id = ?", tournamentID).First(&tournament).Error; err != nil {
		logger.Error(err)
		return false, errors.ErrTournamentNotFound
	}

	return tournament.OwnerId == userID, nil
}

func (tournamentStore *TournamentStore) IsTournamentPlayer(tournamentID uint, userID uint) (bool, error) {
	teams, err := tournamentStore.GetAllTeams(tournamentID)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	for _, team := range *teams {
		for _, player := range team.Players {
			if player.ID == userID {
				return true, nil
			}
		}
	}

	return false, nil
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

	if err = tournamentStore.db.Save(oldTournament).Error; err != nil {
		logger.Error(err)
		return errors.ErrInternal
	}

	return nil
}

func (tournamentStore *TournamentStore) AddTeam(tournamentId uint, teamId uint) error {
	if err := tournamentStore.db.Model(&models.Tournament{ID: tournamentId}).
		Association("Teams").Append(&models.Team{ID: teamId}).Error; err != nil {
		logger.Error(err)
		return errors.ErrTournamentNotFound
	}

	return nil
}

func (tournamentStore *TournamentStore) RemoveTeam(tournamentId uint, teamId uint) error {
	if err := tournamentStore.db.Model(&models.Tournament{ID: tournamentId}).
		Association("Teams").Delete(&models.Team{ID: teamId}).Error; err != nil {
		logger.Error(err)
		return errors.ErrTournamentNotFound
	}

	return nil
}

func (tournamentStore *TournamentStore) GetAllTeams(tournamentId uint) (*models.Teams, error) {
	var tournamentTeams models.Teams
	if err := tournamentStore.db.Model(&models.Tournament{ID: tournamentId}).
		Related(&tournamentTeams, "Teams").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}
	return &tournamentTeams, nil
}

func (tournamentStore *TournamentStore) GetAllMeetings(tournamentId uint) (*models.Meetings, error) {
	var tournamentMeetings models.Meetings
	if err := tournamentStore.db.Model(&models.Tournament{ID: tournamentId}).Preload("Teams").
		Related(&tournamentMeetings, "Meetings").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}

	return &tournamentMeetings, nil
}

func (tournamentStore *TournamentStore) IsTeamInTournament(tournamentId uint, teamId uint) (bool, error) {
	var team models.Team
	if err := tournamentStore.db.Model(&models.Tournament{ID: tournamentId}).
		Where("id = ?", teamId).
		Association("teams").
		Find(&team).Error; err != nil {
		logger.Error(err)
		return false, errors.ErrTeamNotFound
	}
	return true, nil
}
