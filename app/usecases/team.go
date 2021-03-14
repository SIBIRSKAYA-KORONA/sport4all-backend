package usecases

import (
	"sport4all/app/models"
)

type TeamUseCase interface {
	Create(ownerId uint, user *models.Team) error
	GetByID(tid uint) (*models.Team, error)
	GetTeamsByUser(uid uint, role models.Role) (*models.Teams, error)
	GetAllTournaments(tid uint) (*models.Tournaments, error)
	GetTeamsByNamePart(namePart string, limit uint) (*models.Teams, error)
	InviteMember(tid uint, uid uint, role models.Role) error
	GetUsersForInvite(tid uint, nicknamePart string, limit uint) (*models.Users, error)
}
