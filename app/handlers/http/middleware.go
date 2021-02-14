package http

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	useCases "github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/common"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/sanitize"
)

type Middleware interface {
	LogRequest(echo.HandlerFunc) echo.HandlerFunc
	ProcessPanic(echo.HandlerFunc) echo.HandlerFunc
	Sanitize(echo.HandlerFunc) echo.HandlerFunc
	//CORS(next echo.HandlerFunc) echo.HandlerFunс
	CheckAuth(echo.HandlerFunc) echo.HandlerFunc
}

type MiddlewareImpl struct {
	sessionUseCase useCases.SessionUseCase
	//origins    map[string]struct{}
}

func CreateMiddleware(sessionUseCase useCases.SessionUseCase) Middleware {
	return &MiddlewareImpl{sessionUseCase: sessionUseCase}
}

func (mw *MiddlewareImpl) LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
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

func (mw *MiddlewareImpl) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
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

func (mw *MiddlewareImpl) Sanitize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Request().Method != echo.PUT && ctx.Request().Method != echo.POST {
			return next(ctx)
		}
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		defer common.Close(ctx.Request().Body.Close)
		sanBody, err := sanitize.SanitizeJSON(body)
		if err != nil {
			logger.Warn("bluemonday XSS register")
			return ctx.NoContent(http.StatusBadRequest)
		}
		ctx.Set("body", sanBody)
		return next(ctx)
	}
}

func (mw *MiddlewareImpl) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			logger.Error(err)
			return ctx.String(http.StatusUnauthorized, errors.SessionNotFound)
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
