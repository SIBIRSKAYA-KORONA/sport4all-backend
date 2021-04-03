package repositories

import (
	"sport4all/app/models"
)

type TournamentRepository interface {
	Create(tournament *models.Tournament) error
	GetByID(tournamentId uint) (*models.Tournament, error)
	GetTournamentByUserOwner(uid uint) (*models.Tournaments, error)
	Update(tournament *models.Tournament) error
	AddTeam(tournamentId uint, teamId uint) error
	RemoveTeam(tournamentId uint, teamId uint) error
	GetAllTeams(tournamentId uint) (*models.Teams, error)
	GetAllMeetings(tournamentId uint) (*models.Meetings, error)
	IsTournamentOrganizer(tournamentID uint, userID uint) (bool, error)
	IsTournamentPlayer(tournamentID uint, userID uint) (bool, error)
	IsTeamInTournament(tournamentId uint, teamId uint) (bool, error)
}
