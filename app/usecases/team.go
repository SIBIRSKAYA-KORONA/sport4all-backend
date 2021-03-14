package usecases

import (
	"sport4all/app/models"
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
	GetTeamsByUser(uid uint, role Role) (models.Teams, error)
	GetTeamsByNamePart(namePart string, limit uint) (models.Teams, error)
	InviteMember(tid uint, uid uint, role Role) error
	GetUsersForInvite(tid uint, nicknamePart string, limit uint) (models.Users, error)
}
