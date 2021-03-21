package usecases

import (
	"sport4all/app/models"
)

type TournamentUseCase interface {
	Create(tournament *models.Tournament) error
	GetByID(tournamentId uint) (*models.Tournament, error)
	GetTournamentByUser(uid uint) (*models.UserTournament, error)
	Update(meeting *models.Tournament) error
	AddTeam(tournamentId uint, teamId uint) error
	RemoveTeam(tournamentId uint, teamId uint) error
	GetAllTeams(tournamentId uint) (*models.Teams, error)
	GetAllMeetings(tournamentId uint) (*models.Meetings, error)
}
