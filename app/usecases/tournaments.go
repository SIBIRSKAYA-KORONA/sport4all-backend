package usecases

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

const (
	Olympic uint = iota
	Circular
)

type TournamentUseCase interface {
	Create(ownerId uint, tournament *models.Tournament) error
	AddTeam(tournamentId uint, teamId uint) error
	GetByID(tournamentId uint) (*models.Tournament, error)
	GetAllTeams(tournamentId uint) (*models.Teams, error)
	GenerateMeetings(tournamentId uint, genType uint) error
	GetAllMeetings(tournamentId uint) (*models.Meetings, error)
}
