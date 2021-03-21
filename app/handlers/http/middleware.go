package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"sport4all/app/models"
	useCases "sport4all/app/usecases"
	"sport4all/pkg/common"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/sanitize"
)

type Middleware interface {
	LogRequest(echo.HandlerFunc) echo.HandlerFunc
	ProcessPanic(echo.HandlerFunc) echo.HandlerFunc
	Sanitize(echo.HandlerFunc) echo.HandlerFunc
	CORS(echo.HandlerFunc) echo.HandlerFunc
	CheckAuth(echo.HandlerFunc) echo.HandlerFunc
	CheckTeamPermission(role models.Role) echo.MiddlewareFunc
}

type MiddlewareImpl struct {
	sessionUseCase useCases.SessionUseCase
	teamUseCase    useCases.TeamUseCase
	origins        map[string]struct{}
}

func CreateMiddleware(sessionUseCase useCases.SessionUseCase,
	teamUseCase useCases.TeamUseCase,
	origins map[string]struct{}) Middleware {
	return &MiddlewareImpl{
		sessionUseCase: sessionUseCase,
		teamUseCase:    teamUseCase,
		origins:        origins,
	}
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

func (mw *MiddlewareImpl) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		origin := ctx.Request().Header.Get("Origin")
		if _, exist := mw.origins[origin]; !exist {
			return ctx.NoContent(http.StatusForbidden)
		}
		ctx.Response().Header().Set("Access-Control-Allow-Origin", origin)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		if ctx.Request().Method == "OPTIONS" {
			return ctx.NoContent(http.StatusOK)
		}
		return next(ctx)
	}
}

func (mw *MiddlewareImpl) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
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

func (mw *MiddlewareImpl) CheckTeamPermission(role models.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return mw.CheckAuth( func(ctx echo.Context) error {
			var teamID uint
			_, err := fmt.Sscan(ctx.Param("tid"), &teamID)
			if err != nil {
				return ctx.NoContent(http.StatusBadRequest)
			}
			userID := ctx.Get("uid").(uint)

			ok, err := mw.teamUseCase.CheckUserForRole(teamID, userID, role)
			if err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}

			if ok {
				ctx.Set("tid", teamID)
				return next(ctx)
			} else {
				error := errors.ErrNoPermission
				return ctx.String(errors.ResolveErrorToCode(error), error.Error())
			}
		})
	}
}
