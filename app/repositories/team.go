package repositories

import (
	"sport4all/app/models"
)

type TeamRepository interface {
	Create(team *models.Team) error
	GetByID(tid uint) (*models.Team, error)
	GetTeamsByUser(uid uint, role models.Role) (*models.Teams, error)
	GetAllTournaments(tid uint) (*models.Tournaments, error)
	GetTeamsByNamePart(namePart string, limit uint) (*models.Teams, error)
	InviteMember(tid uint, user *models.User, role models.Role) error
	DeleteMember(tid uint, uid uint) error
	GetUsersForInvite(tid uint, nicknamePart string, limit uint) (*models.Users, error)
	IsTeamOwner(teamID uint, userID uint) (bool, error)
	IsTeamPlayer(teamID uint, userID uint) (bool, error)
	GetTeamStats(tid uint) ([]models.Stats, error)
}
