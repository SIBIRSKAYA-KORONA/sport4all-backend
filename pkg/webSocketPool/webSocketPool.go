package webSocketPool

import "github.com/labstack/echo/v4"

type WebSocketPool interface {
	Run(ctx echo.Context) error
	Send(uid uint, mess []byte)
}
