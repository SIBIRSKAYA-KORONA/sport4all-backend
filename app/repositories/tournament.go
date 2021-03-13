package repositories

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type TournamentRepository interface {
	Create(user *models.Tournament) error
	GetByID(tournamentId uint) (*models.Tournament, error)
	AddTeam(tournamentId uint, teamId uint) error
	GetAllTeams(tournamentId uint) (*models.Teams, error)
	GetAllMeetings(tournamentId uint) (*models.Meetings, error)
}
