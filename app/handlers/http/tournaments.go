package http

//
import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/serializer"
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
	tournaments.PUT("/:tournamentId/teams/:tid", handler.AddTeam, mw.CheckAuth)
	tournaments.GET("/:tournamentId/teams", handler.GetAllTeams)
}

func (tournamentHandler *TournamentHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var tournament models.Tournament
	err := serializer.JSON().Unmarshal(body, &tournament)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	ownerId := ctx.Get("uid").(uint)
	err = tournamentHandler.UseCase.Create(ownerId, &tournament)
	if err != nil {
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

func (tournamentHandler *TournamentHandler) AddTeam(ctx echo.Context) error {
	var teamId uint
	_, err := fmt.Sscan(ctx.Param("tid"), &teamId)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var tournamentId uint
	_, err = fmt.Sscan(ctx.Param("tournamentId"), &tournamentId)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	err = tournamentHandler.UseCase.AddTeam(tournamentId, teamId)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return nil
}

func (tournamentHandler *TournamentHandler) GetAllTeams(ctx echo.Context) error {
	return nil
}

func (tournamentHandler *TournamentHandler) GenerateMeetings(ctx echo.Context) error {
	return nil
}

func (tournamentHandler *TournamentHandler) GetAllMeetings(ctx echo.Context) error {
	return nil
}
