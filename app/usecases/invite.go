package usecases

import (
	"sport4all/app/models"
)

type InviteUseCase interface {
	Create(uid uint, invite *models.Invite, entity models.Entity) error
	Update(uid uint, invite *models.Invite) error
	GetUserInvites(uid uint) (*[]models.Invite, error)
	GetTeamInvites(teamId uint, state models.InviteState) (*[]models.Invite, error)
	GetTournamentInvites(tournamentId uint, state models.InviteState) (*[]models.Invite, error)
}
