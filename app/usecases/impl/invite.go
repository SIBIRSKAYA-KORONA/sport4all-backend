package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type InviteUseCaseImpl struct {
	inviteRepo     repositories.InviteRepository
	teamRepo       repositories.TeamRepository
	tournamentRepo repositories.TournamentRepository
}

func CreateInviteUseCase(inviteRepo repositories.InviteRepository,
	teamRepo repositories.TeamRepository,
	tournamentRepo repositories.TournamentRepository) usecases.InviteUseCase {
	return &InviteUseCaseImpl{inviteRepo: inviteRepo, teamRepo: teamRepo, tournamentRepo: tournamentRepo}
}

func (inviteUseCase *InviteUseCaseImpl) Create(uid uint, invite *models.Invite, entity models.Entity) error {
	switch entity {
	case models.TeamEntity:
		return inviteUseCase.createTeamInvite(uid, invite)
	case models.TournamentEntity:
		return inviteUseCase.createTournamentInvite(uid, invite)
	default:
		return nil
	}
}

func (inviteUseCase *InviteUseCaseImpl) createTeamInvite(uid uint, invite *models.Invite) error {
	team, err := inviteUseCase.teamRepo.GetByID((*invite).TeamId)
	if err != nil {
		logger.Error(err)
		return err
	}

	if invite.Type == "direct" {
		isOwner, err := inviteUseCase.teamRepo.IsTeamOwner(team.ID, uid)
		if err != nil {
			logger.Error(err)
			return err
		}

		if !isOwner {
			return errors.ErrNoPermission
		}

		invite.AssignedId = *invite.InvitedId
	} else if invite.Type == "indirect" {
		*invite.InvitedId = uid
		invite.AssignedId = team.OwnerId
	}

	invite.State = models.Opened
	if err := inviteUseCase.inviteRepo.Create(invite); err != nil {
		logger.Info(err)
		return err
	}

	return nil
}

func (inviteUseCase *InviteUseCaseImpl) createTournamentInvite(uid uint, invite *models.Invite) error {
	return nil
}

func (inviteUseCase *InviteUseCaseImpl) Update(uid uint, invite *models.Invite, entity models.Entity) error {
	updatedInvite, err := inviteUseCase.inviteRepo.Update(uid, invite)
	if err != nil {
		logger.Info(err)
		return err
	}

	if updatedInvite != nil {
		if updatedInvite.State == models.Accepted {
			if err := inviteUseCase.teamRepo.InviteMember(updatedInvite.TeamId,
				&models.User{ID: *updatedInvite.InvitedId}, models.Player); err != nil {
				logger.Error(err)
				return err
			}
		}
	}

	return nil
}

func (inviteUseCase *InviteUseCaseImpl) GetUserInvites(uid uint) (*[]models.Invite, error) {
	invites, has := inviteUseCase.inviteRepo.GetUserInvites(uid)
	if !has {
		logger.Info("no messages for the user", uid)
		return nil, nil
	}
	return invites, nil
}

func (inviteUseCase *InviteUseCaseImpl) GetTeamInvites(teamId uint, state models.InviteState) (*[]models.Invite, error) {
	invites, err := inviteUseCase.inviteRepo.GetTeamInvites(teamId, state)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return invites, nil
}
