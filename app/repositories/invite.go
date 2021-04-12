package repositories

import "sport4all/app/models"

type InviteRepository interface {
	Create(invite *models.Invite) error
	Update(uid uint, invite *models.Invite) (*models.Invite, error)
	GetUserInvites(uid uint) (*[]models.Invite, bool)
	GetTeamInvites(teadId uint, state models.InviteState) (*[]models.Invite, error)
}
