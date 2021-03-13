package psql

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
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

func (tournamentStore *TournamentStore) AddTeam(tournamentId uint, teamId uint) error {
	return nil
}

func (tournamentStore *TournamentStore) GetAllTeams(tournamentId uint) (*models.Teams, error) {
	return nil, nil
}

func (tournamentStore *TournamentStore) GetAllMeetings(tournamentId uint) (*models.Meetings, error) {
	return nil, nil
}
