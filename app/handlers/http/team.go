package http

import (
	"fmt"
	"net/http"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/serializer"

	"github.com/labstack/echo/v4"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
)

type TeamHandler struct {
	UseCase  usecases.TeamUseCase
	TeamsURL string
}

func CreateTeamHandler(teamsURL string, router *echo.Group, useCase usecases.TeamUseCase, mw Middleware) {
	handler := &TeamHandler{
		UseCase:  useCase,
		TeamsURL: teamsURL,
	}

	teams := router.Group(handler.TeamsURL)
	teams.POST("", handler.Create, mw.CheckAuth)
	teams.GET("/:tid", handler.GetByID, mw.CheckAuth)

	teams.GET("/search", handler.GetTeamsByNamePart, mw.CheckAuth)

	teams.GET("/:tid/members/search", handler.GetUsersForInvite, mw.CheckAuth)
	teams.POST("/:tid/members/:uid", handler.InviteMember, mw.CheckAuth)
}

func (teamHandler *TeamHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var team models.Team
	err := serializer.JSON().Unmarshal(body, &team)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ownerId := ctx.Get("uid").(uint)
	err = teamHandler.UseCase.Create(ownerId, &team)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&team)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (teamHandler *TeamHandler) GetByID(ctx echo.Context) error {
	var tid uint
	_, err := fmt.Sscan(ctx.Param("tid"), &tid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	team, err := teamHandler.UseCase.GetByID(tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := serializer.JSON().Marshal(&team)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (teamHandler *TeamHandler) GetUsersForInvite(ctx echo.Context) error {
	nicknamePart := ctx.QueryParam("nickname")
	if nicknamePart == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	//bid := ctx.Get("bid").(uint)

	var tid uint
	_, err := fmt.Sscan(ctx.Param("tid"), &tid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var limit uint
	_, err = fmt.Sscan(ctx.QueryParam("limit"), &limit)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	users, err := teamHandler.UseCase.GetUsersForInvite(tid, nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := serializer.JSON().Marshal(&users)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (teamHandler *TeamHandler) InviteMember(ctx echo.Context) error {
	var tid uint
	_, err := fmt.Sscan(ctx.Param("tid"), &tid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var uid uint
	_, err = fmt.Sscan(ctx.Param("uid"), &uid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	roleParam := ctx.QueryParam("role")
	role, exist := usecases.StringToRole[roleParam]
	if !exist {
		return ctx.NoContent(http.StatusBadRequest)
	}

	err = teamHandler.UseCase.InviteMember(tid, uid, role)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (teamHandler *TeamHandler) GetTeamsByNamePart(ctx echo.Context) error {
	namePart := ctx.QueryParam("name")
	if namePart == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}
	var limit uint
	_, err := fmt.Sscan(ctx.QueryParam("limit"), &limit)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	teams, err := teamHandler.UseCase.GetTeamsByNamePart(namePart, limit)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&teams)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
