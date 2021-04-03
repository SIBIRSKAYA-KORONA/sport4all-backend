package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"sport4all/app/usecases"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"
)

type MessageHandler struct {
	UseCase     usecases.MessageUseCase
	MessagesURL string
}

func CreateMessageHandler(messagesURL string, router *echo.Group, useCase usecases.MessageUseCase, mw Middleware) {
	handler := &MessageHandler{
		UseCase:     useCase,
		MessagesURL: messagesURL,
	}

	messages := router.Group(handler.MessagesURL)

	messages.GET("", handler.GetAll, mw.CheckAuth)
	messages.PUT("", handler.UpdateAll, mw.CheckAuth)
	messages.DELETE("", handler.DeleteAll, mw.CheckAuth)
}

func (messageHandler *MessageHandler) GetAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	messages, has := messageHandler.UseCase.GetAll(uid)
	if !has {
		logger.Error("no notifications for the user", uid)
	}
	resp, err := serializer.JSON().Marshal(&messages)
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (messageHandler *MessageHandler) UpdateAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	if err := messageHandler.UseCase.UpdateAll(uid); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (messageHandler *MessageHandler) DeleteAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	if err := messageHandler.UseCase.DeleteAll(uid); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}
