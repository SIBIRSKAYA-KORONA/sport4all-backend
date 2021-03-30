package http

import (
	"fmt"
	"net/http"

	"sport4all/app/models"
	"sport4all/app/usecases"
	"sport4all/pkg/common"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"

	"github.com/labstack/echo/v4"
)

type AttachHandler struct {
	UseCase usecases.AttachUseCase
	URL     string
}

func CreateAttachHandler(attachURL string, router *echo.Group, useCase usecases.AttachUseCase, mw Middleware) {
	handler := &AttachHandler{
		UseCase: useCase,
		URL:     attachURL,
	}

	attach := router.Group(handler.URL)

	attach.PUT("", handler.Create, mw.CheckAuth)
	// attach.DELETE("", handler.Delete, mw.CheckAuth)
}

func getIdFromValue(value string, ctx echo.Context) *uint {
	var id uint
	if _, err := fmt.Sscan(ctx.FormValue(value), &id); err != nil {
		return nil
	}
	return &id
}

func (attachHandler *AttachHandler) Create(ctx echo.Context) error {
	attach := models.Attach{
		UserId:       getIdFromValue("userId", ctx),
		TeamId:       getIdFromValue("teamId", ctx),
		MeetingId:    getIdFromValue("meetingId", ctx),
		TournamentId: getIdFromValue("tournamentId", ctx),
	}

	if attach.UserId == nil && attach.TeamId == nil && attach.MeetingId == nil && attach.TournamentId == nil {
		logger.Error("not set entity id")
		return ctx.String(http.StatusBadRequest, "not set entity id")
	}

	fileHeader, err := ctx.FormFile("attach")
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	attach.Name = fileHeader.Filename
	attach.Data, err = fileHeader.Open()
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	defer common.Close(attach.Data.Close)

	if err = attachHandler.UseCase.Create(&attach); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&attach)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
