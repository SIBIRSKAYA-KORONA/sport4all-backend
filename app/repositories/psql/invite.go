package psql

import (
	"github.com/jinzhu/gorm"
	"sport4all/pkg/errors"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/logger"
)

type InviteStore struct {
	db *gorm.DB
}

func CreateInviteRepository(db *gorm.DB) repositories.InviteRepository {
	return &InviteStore{db: db}
}

func (inviteStore *InviteStore) Create(invite *models.Invite) error {
	if err := inviteStore.db.Create(invite).Error; err != nil {
		logger.Error(err)
		return errors.ErrInternal
	}

	return nil
}

func (inviteStore *InviteStore) GetByID(iid uint) (*models.Invite, error) {
	invite := new(models.Invite)
	if err := inviteStore.db.Where("id = ?", iid).First(&invite).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrInviteNotFound
	}

	return invite, nil
}

func (inviteStore *InviteStore) Update(uid uint, invite *models.Invite) (*models.Invite, error) {
	oldInvite, err := inviteStore.GetByID(invite.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if uid != oldInvite.AssignedId {
		return nil, errors.ErrNoPermission
	}

	newState := invite.State

	if oldInvite.State != newState {
		oldInvite.State = newState

		if err = inviteStore.db.Save(oldInvite).Error; err != nil {
			logger.Error(err)
			return nil, errors.ErrInternal
		}

		return oldInvite, nil
	}

	return nil, nil
}

func (inviteStore *InviteStore) GetUserInvites(uid uint) (*[]models.Invite, bool) {
	var invites []models.Invite
	if err := inviteStore.db.Where("invited_id = ?", uid).Preload("Team.Avatar").Find(&invites).Error; err != nil {
		logger.Error(err)
		return nil, false
	}

	return &invites, true
}

func (inviteStore *InviteStore) GetTeamInvites(teamId uint, state models.InviteState) (*[]models.Invite, error) {
	var invites []models.Invite
	if err := inviteStore.db.Where("team_id = ? AND state = ?", teamId, state).
		Preload("User.Avatar").
		Preload("Tournament.Avatar").
		Find(&invites).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return &invites, nil
}

func (inviteStore *InviteStore) GetTournamentInvites(tournamentId uint, state models.InviteState) (*[]models.Invite, error) {
	var invites []models.Invite
	if err := inviteStore.db.Where("tournament_id = ? AND state = ?", tournamentId, state).
		Preload("Team.Avatar").
		Find(&invites).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return &invites, nil
}
