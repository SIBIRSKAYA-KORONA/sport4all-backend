package usecases

import (
	"sport4all/app/models"
)

type TournamentUseCase interface {
	Create(ownerId uint, tournament *models.Tournament) error
	GetByID(tournamentId uint) (*models.Tournament, error)
	AddTeam(tournamentId uint, teamId uint) error
	GetAllTeams(tournamentId uint) (*models.Teams, error)
	GenerateMeetings(tournamentId uint, genType models.TournamentSystem) error
	GetAllMeetings(tournamentId uint) (*models.Meetings, error)
}
