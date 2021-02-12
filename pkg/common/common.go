package common

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

func Close(closeFunc func() error) {
	if err := closeFunc(); err != nil {
		logger.Error(err)
	}
}

func SetCookie(ctx echo.Context, sid string, expires time.Time) {
	ctx.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
	})
}
