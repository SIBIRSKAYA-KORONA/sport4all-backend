package usecases

import (
	"sport4all/app/models"
)

type TournamentUseCase interface {
	Create(tournament *models.Tournament) error
	GetByID(tournamentId uint) (*models.Tournament, error)
	GetTournamentsByUser(uid uint) (*models.UserTournament, error)
	GetTournamentsByNamePart(namePart string, limit uint) (*models.Tournaments, error)
	Update(meeting *models.Tournament) error
	AddTeam(tournamentId uint, teamId uint) error
	RemoveTeam(tournamentId uint, teamId uint) error
	GetAllTeams(tournamentId uint) (*models.Teams, error)
	GetAllMeetings(tournamentId uint) (*models.Meetings, error)
	CheckUserForTournamentRole(tournamentId uint, uid uint, role models.TournamentRole) (bool, error)
	IsTeamInTournament(tournamentId uint, teamId uint) (bool, error)
	GetTournamentForFeeds(offset uint) (*[]models.Tournament, error)
}
