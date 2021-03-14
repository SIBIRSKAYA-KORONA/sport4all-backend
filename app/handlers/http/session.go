package http

import (
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

type SessionHandler struct {
	UseCase    usecases.SessionUseCase
	SessionURL string
}

func CreateSessionHandler(sessionURL string, router *echo.Group, useCase usecases.SessionUseCase, mw Middleware) {
	handler := &SessionHandler{
		UseCase:    useCase,
		SessionURL: sessionURL,
	}

	settings := router.Group(handler.SessionURL)
	settings.POST("", handler.Create)
	settings.DELETE("", handler.Delete, mw.CheckAuth)
	// TODO: token
}

func (sessionHandler *SessionHandler) Create(ctx echo.Context) error {
	body := ctx.Get("body").([]byte)
	var usr models.User
	if err := serializer.JSON().Unmarshal(body, &usr); err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	session, err := sessionHandler.UseCase.Create(&usr)
	if err != nil {
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	common.SetCookie(ctx, session.SID, time.Now().Add(time.Duration(session.ExpiresSec)*time.Second))
	return ctx.NoContent(http.StatusOK)
}

func (sessionHandler *SessionHandler) Delete(ctx echo.Context) error {
	sid := ctx.Get("sid").(string)
	if err := sessionHandler.UseCase.Delete(sid); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	common.SetCookie(ctx, sid, time.Now().AddDate(-1, 0, 0))
	return ctx.NoContent(http.StatusOK)
}
