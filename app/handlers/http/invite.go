package http

import (
	"fmt"
	"net/http"

	"sport4all/app/models"

	"github.com/labstack/echo/v4"

	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"
)

type InviteHandler struct {
	UseCase   usecases.InviteUseCase
	InviteURL string
}

func CreateInviteHandler(inviteURL string, router *echo.Group, useCase usecases.InviteUseCase, mw Middleware) {
	handler := &InviteHandler{
		UseCase:   useCase,
		InviteURL: inviteURL,
	}

	invites := router.Group(handler.InviteURL)

	invites.POST("/teams", handler.MakeCreateRoute(models.TeamEntity), mw.CheckAuth,
		mw.NotificationMiddleware(models.InviteStatusChanged))
	invites.POST("/tournaments", handler.MakeCreateRoute(models.TournamentEntity), mw.CheckAuth,
		mw.NotificationMiddleware(models.InviteStatusChanged))
	invites.PUT("/:iid", handler.Update, mw.CheckAuth,
		mw.NotificationMiddleware(models.InviteStatusChanged))
	invites.GET("", handler.GetUserInvites, mw.CheckAuth)
	invites.GET("/teams/:tid", handler.GetTeamInvites, mw.CheckAuth)
	invites.GET("/tournaments/:tournamentId", handler.GetTournamentInvites, mw.CheckAuth)
}

func (inviteHandler *InviteHandler) MakeCreateRoute(entity models.Entity) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body := ctx.Get("body").([]byte)
		uid := ctx.Get("uid").(uint)
		var invite models.Invite
		if err := serializer.JSON().Unmarshal(body, &invite); err != nil {
			logger.Error(err)
			return ctx.String(http.StatusBadRequest, err.Error())
		}
		invite.CreatorId = uid
		invite.InviteEntity = entity

		if err := inviteHandler.UseCase.Create(uid, &invite, entity); err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}

		ctx.Set("invite_entity", entity)
		ctx.Set("invite_type", models.InviteType(invite.Type))
		ctx.Set("invite_state", invite.State)
		ctx.Set("assigned", invite.AssignedId)
		ctx.Set("team_id", invite.TeamId)

		if entity == models.TournamentEntity {
			ctx.Set("tournament_id", *invite.TournamentId)
		}
		return ctx.NoContent(http.StatusOK)
	}
}

func (inviteHandler *InviteHandler) Update(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	uid := ctx.Get("uid").(uint)

	var iid uint
	if _, err := fmt.Sscan(ctx.Param("iid"), &iid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var invite models.Invite
	if err := serializer.JSON().Unmarshal(body, &invite); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	invite.ID = iid

	updatedInvite, err := inviteHandler.UseCase.Update(uid, &invite)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	ctx.Set("invite_entity", updatedInvite.InviteEntity)
	ctx.Set("invite_type", models.InviteType(updatedInvite.Type))
	ctx.Set("invite_state", updatedInvite.State)
	ctx.Set("assigned", updatedInvite.AssignedId)
	ctx.Set("author", updatedInvite.CreatorId)
	ctx.Set("team_id", updatedInvite.TeamId)

	if updatedInvite.InviteEntity == models.TournamentEntity {
		ctx.Set("tournament_id", *updatedInvite.TournamentId)
	}

	return ctx.NoContent(http.StatusOK)
}

func (inviteHandler *InviteHandler) GetUserInvites(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	invites, err := inviteHandler.UseCase.GetUserInvites(uid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&invites)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

var stringToState = map[string]models.InviteState{
	"opened":   models.Opened,
	"rejected": models.Rejected,
	"accepted": models.Accepted,
}

func (inviteHandler *InviteHandler) GetTeamInvites(ctx echo.Context) error {
	stateStr := ctx.QueryParam("state")
	state, has := stringToState[stateStr]
	if !has {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var teamId uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &teamId); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	invites, err := inviteHandler.UseCase.GetTeamInvites(teamId, state)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&invites)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (inviteHandler *InviteHandler) GetTournamentInvites(ctx echo.Context) error {
	stateStr := ctx.QueryParam("state")
	state, has := stringToState[stateStr]
	if !has {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var tournamentId uint
	if _, err := fmt.Sscan(ctx.Param("tournamentId"), &tournamentId); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	invites, err := inviteHandler.UseCase.GetTournamentInvites(tournamentId, state)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&invites)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
