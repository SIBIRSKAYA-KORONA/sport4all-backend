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

func (tournamentStore *TournamentStore) AddTeam(tournamentId uint, teamId uint) error {
	team := new(models.Team)
	if err := tournamentStore.DB.First(team, teamId).Error; err != nil {
		logger.Error(err)
		return errors.ErrTeamNotFound
	}

	if err := tournamentStore.DB.Model(&models.Tournament{ID: tournamentId}).
		Association("Teams").Append(team).Error; err != nil { // TODO: mey be Append(&models.Team{ID: teamId})
		logger.Error(err)
		return errors.ErrInternal
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
	if err := tournamentStore.DB.Model(&models.Tournament{ID: tournamentId}).
		Related(&tournamentMeetings, "Meetings").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTournamentNotFound
	}
	return &tournamentMeetings, nil
}
