package ws

import (
	"sync"

	"github.com/labstack/echo/v4"
	"sport4all/app/handlers/http"
	"sport4all/pkg/logger"
	"sport4all/pkg/webSocketPool"
)

type Api interface {
	Run(wg *sync.WaitGroup)
}

type ApiImpl struct {
	wsPool webSocketPool.WebSocketPool
	mw     http.MiddlewareMini
}

func CreateApi(wsPool webSocketPool.WebSocketPool, mw http.MiddlewareMini) Api {
	return &ApiImpl{wsPool: wsPool, mw: mw}
}

func (api *ApiImpl) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	router := echo.New()

	router.Use(api.mw.ProcessPanic)
	router.Use(api.mw.LogRequest)

	router.GET("/api/ws", api.Handle, api.mw.CheckAuth)

	if err := router.Start("0.0.0.0:6060"); err != nil {
		logger.Fatal(err)
	}
}

func (api *ApiImpl) Handle(ctx echo.Context) error {
	return api.wsPool.Run(ctx)
}
