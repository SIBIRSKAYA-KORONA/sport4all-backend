package usecases

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type Role uint8

const (
	Player Role = iota
	Owner
)

var StringToRole = map[string]Role{
	"player": Player,
	"owner":  Owner,
}

type TeamUseCase interface {
	Create(ownerId uint, user *models.Team) error
	GetByID(tid uint) (*models.Team, error)
	GetTeamsByNamePart(namePart string, limit uint) (models.Teams, error)
	InviteMember(tid uint, uid uint, role Role) error
	GetUsersForInvite(tid uint, nicknamePart string, limit uint) (models.Users, error)
}
