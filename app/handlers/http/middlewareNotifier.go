package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	useCases "sport4all/app/usecases"
	"sport4all/pkg/common"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type MiddlewareMini interface {
	LogRequest(echo.HandlerFunc) echo.HandlerFunc
	ProcessPanic(echo.HandlerFunc) echo.HandlerFunc
	CheckAuth(echo.HandlerFunc) echo.HandlerFunc
}

type MiddlewareMiniImpl struct {
	sessionUseCase useCases.SessionUseCase
}

func CreateMiddlewareMini(sessionUseCase useCases.SessionUseCase) MiddlewareMini {
	return &MiddlewareMiniImpl{
		sessionUseCase: sessionUseCase,
	}
}

func (mw *MiddlewareMiniImpl) LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		res := next(ctx)
		logger.Infof("%s %s %d %s",
			ctx.Request().Method,
			ctx.Request().RequestURI,
			ctx.Response().Status,
			time.Since(start))
		return res
	}
}

func (mw *MiddlewareMiniImpl) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Process panic up on: ", ctx.Request().Method,
					ctx.Request().URL.Path, " statement: ", err)
				if err = ctx.NoContent(http.StatusInternalServerError); err != nil {
					logger.Error(err)
				}
			}
		}()
		return next(ctx)
	}
}

func (mw *MiddlewareMiniImpl) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(errors.ErrSessionNotFound), errors.ErrSessionNotFound.Error())
		}
		sid := cookie.Value
		uid, err := mw.sessionUseCase.GetByID(sid)
		if err != nil {
			logger.Error(err)
			common.SetCookie(ctx, sid, time.Now().AddDate(-1, 0, 0))
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		ctx.Set("uid", uid)
		ctx.Set("sid", sid)
		return next(ctx)
	}
}
