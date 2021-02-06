package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/middleware"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/serializer"
)

type UserHandler struct {
	UseCase usecases.UserUseCase
	BaseURL string
}

func CreateUserHandler(group *echo.Group, baseURL string, useCase usecases.UserUseCase, mw middleware.Middleware) {
	handler := &UserHandler{
		UseCase: useCase,
		BaseURL: baseURL,
	}

	// Группа
	settings := group.Group(handler.BaseURL)
	// Общие для всей группы миддлвары
	settings.Use(mw.CheckAuth)

	// Регистрация endpoint'ов в группе

	// Прямо в корень
	settings.POST("", handler.Create, mw.Sanitize)
	settings.PUT("", handler.Create, mw.Sanitize)
	settings.GET("", handler.Create, mw.Sanitize)

	// С расширением урла
	settings.DELETE("/delete", handler.Create, mw.Sanitize)

	group.GET("/profile/:id_or_nickname", handler.Get)
}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	var usr models.User
	body := ctx.Get("body").([]byte)

	err := serializer.JSON().Unmarshal(body, &usr)
	logger.Debug(usr)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	err = userHandler.UseCase.Create(&usr)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Get(ctx echo.Context) error {
	usrKey := ctx.Param("id_or_nickname")
	var usr *models.User
	var err error
	if _, er := strconv.Atoi(usrKey); er == nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		usr, err = userHandler.UseCase.GetByNickname(usrKey)
	}
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := serializer.JSON().Marshal(&usr)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
