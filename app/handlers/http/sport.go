package http

import (
	"net/http"

	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"

	"github.com/labstack/echo/v4"
)

type SportHandler struct {
	UseCase usecases.SportUseCase
	URL     string
}

func CreateSportHandler(url string, router *echo.Group, useCase usecases.SportUseCase, mw Middleware) {
	handler := &SportHandler{
		UseCase: useCase,
		URL:     url,
	}

	attach := router.Group(handler.URL)

	attach.GET("", handler.GetAll)
	// attach.DELETE("", handler.Delete, mw.CheckAuth)
}

func (sportHandler *SportHandler) Create(ctx echo.Context) error {
	return ctx.String(http.StatusInternalServerError, "not impl method")
}

func (sportHandler *SportHandler) GetAll(ctx echo.Context) error {
	sports, err := sportHandler.UseCase.GetAll(50)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&sports)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
