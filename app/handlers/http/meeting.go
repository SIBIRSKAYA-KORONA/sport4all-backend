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

	// --- CRUD ---
	meeting.POST("", handler.Create, mw.CheckAuth)
	meeting.GET("/:mid", handler.GetByID)
	meeting.PUT("/:mid", handler.Update, mw.CheckTournamentPermissionByMeeting(models.TournamentOrganizer),
		mw.NotificationMiddleware(models.EventStatusChanged))

	// --- Управление командами во встрече ---
	meeting.POST("/:mid/teams/:tid", handler.AssignTeam, mw.CheckMeetingStatus(models.RegistrationEvent))

	// --- Статистика ---
	meeting.GET("/:mid/stat", handler.GetMeetingStat)
	meeting.PUT("/:mid/teams/:tid/stat", handler.CreateTeamStat,
		mw.CheckMeetingStatus(models.InProgressEvent), mw.CheckTeamInMeeting)
	meeting.PUT("/:mid/teams/:tid/players/:uid/stat", handler.CreatePlayerStat,
		mw.CheckMeetingStatus(models.InProgressEvent), mw.CheckPlayerInTeam())

	meeting.POST("/:mid/players/stats", handler.CreatePlayersStats, mw.CheckMeetingStatus(models.InProgressEvent))
	meeting.GET("/:mid/players/stats", handler.GetStatsByImage, mw.CheckMeetingStatus(models.InProgressEvent))
}

func (meetingHandler *MeetingHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var meeting models.Meeting
	if err := serializer.JSON().Unmarshal(body, &meeting); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if err := meetingHandler.UseCase.Create(&meeting); err != nil {
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
	if _, err := fmt.Sscan(ctx.Param("mid"), &mid); err != nil {
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
	body := ctx.Get("body").([]byte)
	var meeting models.Meeting
	if err := serializer.JSON().Unmarshal(body, &meeting); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	if _, err := fmt.Sscan(ctx.Param("mid"), &meeting.ID); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	ctx.Set("status", uint(meeting.Status))
	ctx.Set("event_entity", models.MeetingEntity)

	if err := meetingHandler.UseCase.Update(&meeting); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (meetingHandler *MeetingHandler) AssignTeam(ctx echo.Context) error {
	mid := ctx.Get("meetingId").(uint)

	var tid uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &tid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err := meetingHandler.UseCase.AssignTeam(mid, tid); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (meetingHandler *MeetingHandler) CreateTeamStat(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var stats models.Stats
	if err := serializer.JSON().Unmarshal(body, &stats); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	stats.TeamId = ctx.Get("teamId").(uint)
	stats.MeetingId = ctx.Get("meetingId").(uint)

	if err := meetingHandler.UseCase.CreateTeamStat(&stats); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (meetingHandler *MeetingHandler) CreatePlayerStat(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var stats models.Stats
	if err := serializer.JSON().Unmarshal(body, &stats); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	stats.TeamId = ctx.Get("teamId").(uint)
	stats.MeetingId = ctx.Get("meetingId").(uint)
	uid := ctx.Get("playerId").(uint)
	stats.PlayerId = &uid

	if err := meetingHandler.UseCase.CreateTeamStat(&stats); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (meetingHandler *MeetingHandler) CreatePlayersStats(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var stats []models.Stats
	if err := serializer.JSON().Unmarshal(body, &stats); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	mid := ctx.Get("meetingId").(uint)
	if err := meetingHandler.UseCase.CreatePlayersStats(mid, &stats); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&stats)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (meetingHandler *MeetingHandler) GetMeetingStat(ctx echo.Context) error {
	var mid uint
	if _, err := fmt.Sscan(ctx.Param("mid"), &mid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	stats, err := meetingHandler.UseCase.GetMeetingStat(mid)
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

func (meetingHandler *MeetingHandler) GetStatsByImage(ctx echo.Context) error {
	protocolType := ctx.QueryParam("protocol")
	imagePath := ctx.QueryParam("path")
	meetingId := ctx.Get("meetingId").(uint)
	stats, err := meetingHandler.UseCase.GetStatsByImage(meetingId, imagePath, protocolType)
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
