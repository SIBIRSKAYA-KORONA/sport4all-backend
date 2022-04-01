package http

import (
	"fmt"
	"net/http"

	"sport4all/app/models"
	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"

	"github.com/labstack/echo/v4"
)

type SkillHandler struct {
	UseCase usecases.SkillUseCase
	URL     string
}

func CreateSkillHandler(url string, router *echo.Group, useCase usecases.SkillUseCase, mw Middleware) {
	handler := &SkillHandler{
		UseCase: useCase,
		URL:     url,
	}

	group := router.Group(handler.URL)

	group.POST("/:uid", handler.Create, mw.CheckAuth)
	group.GET("/search", handler.GetByNamePart)
	group.POST("/:sid/approve/:uid", handler.CreateApprove, mw.CheckAuth, mw.NotificationMiddleware(models.SkillApproved))
	// group.DELETE("/:sid", handler.Delete, mw.CheckAuth)
	// group.DELETE("/:sid/approve", handler.DeleteApprove, mw.CheckAuth)
}

func (skillHandler *SkillHandler) Create(ctx echo.Context) error {
	var toUid uint
	if _, err := fmt.Sscan(ctx.Param("uid"), &toUid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	body := ctx.Get("body").([]byte)
	var skill models.Skill
	if err := serializer.JSON().Unmarshal(body, &skill); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	fromUid := ctx.Get("uid").(uint)
	if err := skillHandler.UseCase.Create(toUid, fromUid, &skill); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&skill)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (skillHandler *SkillHandler) GetByNamePart(ctx echo.Context) error {
	namePart := ctx.QueryParam("name")
	if namePart == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var limit uint
	if _, err := fmt.Sscan(ctx.QueryParam("limit"), &limit); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	skills, err := skillHandler.UseCase.GetByNamePart(namePart, limit)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&skills)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (skillHandler *SkillHandler) Delete(ctx echo.Context) error {
	return ctx.NoContent(http.StatusInternalServerError)
}

func (skillHandler *SkillHandler) CreateApprove(ctx echo.Context) error {
	var sid uint
	if _, err := fmt.Sscan(ctx.Param("sid"), &sid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	body := ctx.Get("body").([]byte)
	var approve models.SkillApprove
	if err := serializer.JSON().Unmarshal(body, &approve); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if _, err := fmt.Sscan(ctx.Param("uid"), &approve.ToUid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	approve.SkillId = sid
	approve.FromUid = ctx.Get("uid").(uint)

	if err := skillHandler.UseCase.CreateApprove(&approve); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	ctx.Set("toUid", approve.ToUid)

	resp, err := serializer.JSON().Marshal(&approve)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (skillHandler *SkillHandler) DeleteApprove(ctx echo.Context) error {
	return ctx.NoContent(http.StatusInternalServerError)
}
