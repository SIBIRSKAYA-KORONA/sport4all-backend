package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"sport4all/app/models"
	"sport4all/app/usecases"
	"sport4all/pkg/common"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"
)

type UserHandler struct {
	UseCase     usecases.UserUseCase
	SettingsURL string
	ProfileURL  string
}

func CreateUserHandler(settingsURL string, profileURL string, router *echo.Group, useCase usecases.UserUseCase, mw Middleware) {
	handler := &UserHandler{
		UseCase:     useCase,
		SettingsURL: settingsURL,
		ProfileURL:  profileURL,
	}

	router.GET(handler.ProfileURL+"/:nickname", handler.GetByNickname)

	settings := router.Group(handler.SettingsURL)
	settings.POST("", handler.Create)
	settings.GET("", handler.GetByID, mw.CheckAuth)
	settings.PUT("", handler.Update, mw.CheckAuth)
	settings.DELETE("", handler.Delete, mw.CheckAuth)

	// --- Статистика ---
	settings.GET("/:uid/stats", handler.GetUserStats, mw.CheckAuth)
}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var usr models.User
	if err := serializer.JSON().Unmarshal(body, &usr); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	session, err := userHandler.UseCase.Create(&usr)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	common.SetCookie(ctx, session.SID, time.Now().Add(time.Duration(session.ExpiresSec)*time.Second))
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) GetByNickname(ctx echo.Context) error {
	usrKey := ctx.Param("nickname")
	usr, err := userHandler.UseCase.GetByNickname(usrKey)
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

func (userHandler *UserHandler) GetByID(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	usr, err := userHandler.UseCase.GetByID(uid)
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

func (userHandler *UserHandler) Update(ctx echo.Context) error {
	// TODO:
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Delete(ctx echo.Context) error {
	// TODO:
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) GetUserStats(ctx echo.Context) error {
	var uid uint
	if _, err := fmt.Sscan(ctx.Param("uid"), &uid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	stats, err := userHandler.UseCase.GetUserStats(uid)
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
