package repositories

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
)

type TeamRepository interface {
	Create(user *models.Team) error
	GetByID(tid uint) (*models.Team, error)
	GetTeamsByNamePart(namePart string, limit uint) (models.Teams, error)
	InviteMember(tid uint, user *models.User, role usecases.Role) error
}
