package http

import (
	"net/http"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/serializer"

	"github.com/labstack/echo/v4"
)

type MeetingStatHandler struct {
	UseCase         usecases.MeetingStatUseCase
	MeetingsURL     string
	MeetingsStatURL string
}

func CreateMeetingStatHandler(meetingsURL string, meetingsStatURL string, router *echo.Group, useCase usecases.MeetingStatUseCase, mw Middleware) {
	handler := &MeetingStatHandler{
		UseCase:         useCase,
		MeetingsURL:     meetingsURL,
		MeetingsStatURL: meetingsStatURL,
	}

	meeting := router.Group(handler.MeetingsURL)

	meeting.POST("/:mid"+meetingsStatURL, handler.Create, mw.CheckAuth, mw.CheckMeeting)
	//meeting.GET("/:mid" + , handler.GetByID)
	//meeting.PUT("/:mid", handler.Update, mw.CheckAuth)

}

func (meetingStatHandler *MeetingStatHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var meetingStat models.MeetingStat
	err := serializer.JSON().Unmarshal(body, &meetingStat)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	meetingID := ctx.Get("meetingID").(uint)
	meetingStat.MeetingId = meetingID

	err = meetingStatHandler.UseCase.Create(&meetingStat)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&meetingStat)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

//func (meetingStatHandler *MeetingStatHandler) GetByID(ctx echo.Context) error {
//	var mid uint
//	if _, err := fmt.Sscan(ctx.Param("tid"), &mid); err != nil {
//		return ctx.NoContent(http.StatusBadRequest)
//	}
//
//	meeting, err := meetingStatHandler.UseCase.GetByID(mid)
//	if err != nil {
//		logger.Error(err)
//		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
//	}
//	resp, err := serializer.JSON().Marshal(&meeting)
//	if err != nil {
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//	return ctx.String(http.StatusOK, string(resp))
//}
//
//func (meetingStatHandler *MeetingStatHandler) Update(ctx echo.Context) error {
//	return ctx.String(http.StatusOK, "")
//}
