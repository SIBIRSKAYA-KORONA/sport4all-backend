package http

import (
	"net"
	"net/http"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"
	"strings"

	"github.com/labstack/echo/v4"

	"sport4all/app/models"
	"sport4all/app/usecases"
	// "sport4all/pkg/errors"
	// "sport4all/pkg/logger"
	// "sport4all/pkg/serializer"
)

type SearchHandler struct {
	UseCase   usecases.SearchUseCase
	SearchURL string
}

func CreateSearchHandler(searchURL string, router *echo.Group, useCase usecases.SearchUseCase, enableGeo bool, mw Middleware) {
	handler := &SearchHandler{
		UseCase:   useCase,
		SearchURL: searchURL,
	}

	search := router.Group(handler.SearchURL)
	search.GET("", handler.GetResult)
	if enableGeo {
		search.GET("/geo", handler.GetGeo)
	}
}

func (searchHandler *SearchHandler) GetGeo(ctx echo.Context) error {
	ctxIP := ctx.QueryParam("ip")
	if ctxIP == "" {
		ctxIP = ctx.RealIP()
	}
	ip := net.ParseIP(ctxIP)
	if ip == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	
	location, err := searchHandler.UseCase.GetGeo(ip)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&location)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (searchHandler *SearchHandler) GetResult(ctx echo.Context) error {
	var uid *uint
	if value, ok := ctx.Get("uid").(uint); ok {
		*uid = value
	} else {
		uid = nil
	}
	entitiesParam := ctx.QueryParam("entities")
	entitiesMap := searchHandler.parseEntities(entitiesParam)
	if len(entitiesMap) == 0 {
		return ctx.NoContent(http.StatusBadRequest)
	}
	text := ctx.QueryParam("text")
	if text == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	baseOfQuery := &models.SearchQueryBase{Text: text, Offset: 10}

	searchOutput, err := searchHandler.UseCase.GetResult(uid, searchHandler.processEntities(ctx, entitiesMap, baseOfQuery))

	resp, err := serializer.JSON().Marshal(&searchOutput)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (searchHandler *SearchHandler) parseEntities(str string) map[models.Entity]bool {
	str += "|"
	splited := strings.Split(str, "|")
	if len(splited) == 0 {
		return nil
	}

	entityMap := map[models.Entity]bool{}

	for index, _ := range splited {
		if entity, ok := models.StrToEntity[splited[index]]; ok {
			entityMap[entity] = true
		}
	}

	return entityMap
}

func (searchHandler *SearchHandler) processEntities(ctx echo.Context, entities map[models.Entity]bool, base *models.SearchQueryBase) *models.SearchInput {
	input := new(models.SearchInput)
	if _, ok := entities[models.TeamEntity]; ok {
		input.TeamQuery = searchHandler.parseTeamQuery(ctx, base)
	}
	if _, ok := entities[models.TournamentEntity]; ok {
		input.TournamentQuery = searchHandler.parseTournamentQuery(ctx, base)
	}
	if _, ok := entities[models.UserEntity]; ok {
		input.UserQuery = searchHandler.parseUserQuery(ctx, base)
	}
	return input
}

func (searchHandler *SearchHandler) parseUserQuery(ctx echo.Context, base *models.SearchQueryBase) *models.UserSearchQuery {
	return &models.UserSearchQuery{
		Base: base,
		// Role
	}
}
func (searchHandler *SearchHandler) parseTeamQuery(ctx echo.Context, base *models.SearchQueryBase) *models.TeamSearchQuery {
	return &models.TeamSearchQuery{
		Base: base,
		// Location
	}
}
func (searchHandler *SearchHandler) parseTournamentQuery(ctx echo.Context, base *models.SearchQueryBase) *models.TournamentSearchQuery {
	return &models.TournamentSearchQuery{
		Base:        base,
		KindOfSport: ctx.QueryParam("sportKind"),
		// Location
	}
}
