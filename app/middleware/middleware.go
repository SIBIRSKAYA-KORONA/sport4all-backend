package middleware

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

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

type implementation struct {
	// CORS
	//origins    map[string]struct{}
}

func CreateMiddleware() Middleware {
	//origins_ := make(map[string]struct{})
	//for _, key := range viper.GetStringSlice("cors.allowed_origins") {
	//	origins_[key] = struct{}{}
	//}

	return &implementation{}
}

func (mw *implementation) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		logger.Debug("CheckAuth called")
		return next(ctx)
	}
}

func (mw *implementation) LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
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

func (mw *implementation) ProcessPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("ProcessPanic up on ", ctx.Request().Method, ctx.Request().URL.Path)
				logger.Error("Panic statement: ", err)
				err := ctx.NoContent(http.StatusInternalServerError)
				if err != nil {
					logger.Error(err)
				}
			}
		}()
		return next(ctx)
	}
}

func (mw *implementation) Sanitize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		defer ctx.Request().Body.Close()
		sanBody, err := sanitize.SanitizeJSON(body)
		if err != nil {
			logger.Warn("bluemonday XSS register")
			return ctx.NoContent(http.StatusBadRequest)
		}
		ctx.Set("body", sanBody)
		return next(ctx)
	}
}
