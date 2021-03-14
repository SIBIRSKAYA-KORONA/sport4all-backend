package repositories

import (
	"sport4all/app/models"
	"sport4all/app/usecases"
)

type TeamRepository interface {
	Create(team *models.Team) error
	GetByID(tid uint) (*models.Team, error)
	GetTeamsByUser(uid uint, role usecases.Role) (models.Teams, error)
	GetTeamsByNamePart(namePart string, limit uint) (models.Teams, error)
	InviteMember(tid uint, user *models.User, role usecases.Role) error
	GetUsersForInvite(tid uint, nicknamePart string, limit uint) (models.Users, error)
}
