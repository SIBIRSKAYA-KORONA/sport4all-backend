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

type MeetingHandler struct {
	UseCase     usecases.MeetingUseCase
	MeetingsURL string
}

func CreateMeetingsHandler(meetingsURL string, router *echo.Group, useCase usecases.MeetingUseCase, mw Middleware) {
	handler := &MeetingHandler{
		UseCase:     useCase,
		MeetingsURL: meetingsURL,
	}

	meeting := router.Group(handler.MeetingsURL)

	meeting.POST("", handler.Create, mw.CheckAuth)
	meeting.GET("/:mid", handler.GetByID)
	meeting.PUT("/:mid", handler.Update, mw.CheckAuth)
}

func (meetingHandler *MeetingHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var meeting models.Meeting
	err := serializer.JSON().Unmarshal(body, &meeting)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	err = meetingHandler.UseCase.Create(&meeting)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&meeting)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (meetingHandler *MeetingHandler) GetByID(ctx echo.Context) error {
	var mid uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &mid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	meeting, err := meetingHandler.UseCase.GetByID(mid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := serializer.JSON().Marshal(&meeting)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (meetingHandler *MeetingHandler) Update(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "")
}
