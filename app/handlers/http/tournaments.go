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

type TournamentHandler struct {
	UseCase        usecases.TournamentUseCase
	TournamentsURL string
}

func CreateTournamentHandler(tournamentsURL string, router *echo.Group, useCase usecases.TournamentUseCase, mw Middleware) {
	handler := &TournamentHandler{
		UseCase:        useCase,
		TournamentsURL: tournamentsURL,
	}

	tournaments := router.Group(handler.TournamentsURL)
	tournaments.POST("", handler.Create, mw.CheckAuth)
	tournaments.GET("", handler.GetTournamentByUser)
	tournaments.GET("/:tournamentId", handler.GetByID)
	tournaments.PUT("/:tournamentId", handler.Update, mw.CheckAuth)
	tournaments.PUT("/:tournamentId/teams/:tid", handler.AddTeam, mw.CheckTournamentPermission(models.TournamentOrganizer))
	tournaments.DELETE("/:tournamentId/teams/:tid", handler.RemoveTeam, mw.CheckTournamentPermission(models.TournamentOrganizer))
	tournaments.GET("/:tournamentId/teams", handler.GetAllTeams)
	tournaments.GET("/:tournamentId/meetings", handler.GetAllMeetings)
	tournaments.GET("/feeds", handler.GetTournamentForFeeds)
}

func (tournamentHandler *TournamentHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var tournament models.Tournament
	if err := serializer.JSON().Unmarshal(body, &tournament); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	tournament.OwnerId = ctx.Get("uid").(uint)

	if err := tournamentHandler.UseCase.Create(&tournament); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&tournament)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (tournamentHandler *TournamentHandler) GetTournamentByUser(ctx echo.Context) error {
	var userId uint
	if _, err := fmt.Sscan(ctx.QueryParam("userId"), &userId); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	userTournament, err := tournamentHandler.UseCase.GetTournamentByUser(userId)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := serializer.JSON().Marshal(&userTournament)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (tournamentHandler *TournamentHandler) GetByID(ctx echo.Context) error {
	var tournamentId uint
	if _, err := fmt.Sscan(ctx.Param("tournamentId"), &tournamentId); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	tournament, err := tournamentHandler.UseCase.GetByID(tournamentId)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := serializer.JSON().Marshal(&tournament)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (tournamentHandler *TournamentHandler) Update(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var tournament models.Tournament
	if err := serializer.JSON().Unmarshal(body, &tournament); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if _, err := fmt.Sscan(ctx.Param("tournamentId"), &tournament.ID); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err := tournamentHandler.UseCase.Update(&tournament); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (tournamentHandler *TournamentHandler) AddTeam(ctx echo.Context) error {
	var teamId uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &teamId); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	tournamentId := ctx.Get("tournamentId").(uint)

	if err := tournamentHandler.UseCase.AddTeam(tournamentId, teamId); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (tournamentHandler *TournamentHandler) RemoveTeam(ctx echo.Context) error {
	var teamId uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &teamId); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	tournamentId := ctx.Get("tournamentId").(uint)

	if err := tournamentHandler.UseCase.RemoveTeam(tournamentId, teamId); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (tournamentHandler *TournamentHandler) GetAllTeams(ctx echo.Context) error {
	var tid uint
	if _, err := fmt.Sscan(ctx.Param("tournamentId"), &tid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	teams, err := tournamentHandler.UseCase.GetAllTeams(tid)
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

func (tournamentHandler *TournamentHandler) GetAllMeetings(ctx echo.Context) error {
	var tid uint
	if _, err := fmt.Sscan(ctx.Param("tournamentId"), &tid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	meetings, err := tournamentHandler.UseCase.GetAllMeetings(tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&meetings)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (tournamentHandler *TournamentHandler) GetTournamentForFeeds(ctx echo.Context) error {
	var offset uint
	if _, err := fmt.Sscan(ctx.QueryParam("offset"), &offset); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	tournaments, err := tournamentHandler.UseCase.GetTournamentForFeeds(offset, 10)
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
