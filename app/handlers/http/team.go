package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"sport4all/app/models"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"
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

	// --- CRUD ---
	teams.POST("", handler.Create, mw.CheckAuth)
	teams.GET("", handler.GetTeamsByUser, mw.CheckAuth)
	teams.GET("/:tid", handler.GetByID)
	teams.GET("/:tid/tournaments", handler.GetAllTournaments)

	// --- Поиск ---
	teams.GET("/search", handler.GetTeamsByNamePart)
	teams.GET("/:tid/members/search", handler.GetUsersForInvite, mw.CheckTeamPermission(models.Owner))

	// --- Управление участниками ---
	teams.POST("/:tid/members/:uid", handler.InviteMember, mw.CheckTeamPermission(models.Owner), mw.NotificationMiddleware)
	teams.DELETE("/:tid/members/:uid", handler.DeleteMember, mw.CheckTeamPermission(models.Owner))

	// --- Статистика ---
	teams.GET("/:tid/stats", handler.GetTeamStats, mw.CheckAuth)
}

func (teamHandler *TeamHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var team models.Team
	if err := serializer.JSON().Unmarshal(body, &team); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	ownerId := ctx.Get("uid").(uint)
	if err := teamHandler.UseCase.Create(ownerId, &team); err != nil {
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

func (teamHandler *TeamHandler) GetTeamsByUser(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	role, exist := models.StringToRole[ctx.QueryParam("role")]
	if !exist {
		return ctx.NoContent(http.StatusBadRequest)
	}

	teams, err := teamHandler.UseCase.GetTeamsByUser(uid, role)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&teams)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (teamHandler *TeamHandler) GetByID(ctx echo.Context) error {
	var tid uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &tid); err != nil {
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

func (teamHandler *TeamHandler) GetAllTournaments(ctx echo.Context) error {
	var tid uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &tid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	tournaments, err := teamHandler.UseCase.GetAllTournaments(tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&tournaments)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (teamHandler *TeamHandler) GetUsersForInvite(ctx echo.Context) error {
	nicknamePart := ctx.QueryParam("nickname")
	if nicknamePart == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	tid := ctx.Get("tid").(uint)

	var limit uint
	if _, err := fmt.Sscan(ctx.QueryParam("limit"), &limit); err != nil {
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
	tid := ctx.Get("tid").(uint)

	var uid uint
	if _, err := fmt.Sscan(ctx.Param("uid"), &uid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	roleParam := ctx.QueryParam("role")
	role, exist := models.StringToRole[roleParam]
	if !exist {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err := teamHandler.UseCase.InviteMember(tid, uid, role); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	ctx.Set("member", uid)
	return ctx.NoContent(http.StatusOK)
}

func (teamHandler *TeamHandler) DeleteMember(ctx echo.Context) error {
	tid := ctx.Get("tid").(uint)

	var uid uint
	if _, err := fmt.Sscan(ctx.Param("uid"), &uid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err := teamHandler.UseCase.DeleteMember(tid, uid); err != nil {
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
	if _, err := fmt.Sscan(ctx.QueryParam("limit"), &limit); err != nil {
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

func (teamHandler *TeamHandler) GetTeamStats(ctx echo.Context) error {
	var tid uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &tid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	stats, err := teamHandler.UseCase.GetTeamStats(tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&stats)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
