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
	UseCase        usecases.UserUseCase
	SessionUseCase usecases.SessionUseCase
	SettingsURL    string
	ProfileURL     string
}

func CreateUserHandler(settingsURL string, profileURL string, router *echo.Group, useCase usecases.UserUseCase,
	sessionUseCase usecases.SessionUseCase,
	mw Middleware) {
	handler := &UserHandler{
		UseCase:        useCase,
		SessionUseCase: sessionUseCase,
		SettingsURL:    settingsURL,
		ProfileURL:     profileURL,
	}

	profile := router.Group(handler.ProfileURL)
	profile.GET("/:nickname", handler.GetByNickname)
	profile.GET("/:uid/skills", handler.GetUserSkills)
	profile.GET("/:uid/stats", handler.GetUserStats)
	profile.GET("/search", handler.SearchUsers)

	settings := router.Group(handler.SettingsURL)
	settings.POST("", handler.Create)
	settings.GET("", handler.GetByID, mw.CheckAuth)
	settings.PUT("", handler.Update, mw.CheckAuth)
	settings.DELETE("", handler.Delete, mw.CheckAuth)
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
	var newUser models.User

	newUser.ID = ctx.Get("uid").(uint)
	newUser.Name = ctx.FormValue("newName")
	newUser.Surname = ctx.FormValue("newSurname")
	newUser.About = ctx.FormValue("about")
	newUser.Birthday = ctx.FormValue("birthday")
	newUser.PhoneNumber = ctx.FormValue("phoneNumber")

	if err := userHandler.UseCase.Update(&newUser); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) Delete(ctx echo.Context) error {
	// TODO:
	return ctx.NoContent(http.StatusOK)
}

func (userHandler *UserHandler) GetUserSkills(ctx echo.Context) error {
	var uid uint
	if _, err := fmt.Sscan(ctx.Param("uid"), &uid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	skills, err := userHandler.UseCase.GetUserSkills(uid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	resp, err := serializer.JSON().Marshal(&skills)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
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

func (userHandler *UserHandler) SearchUsers(ctx echo.Context) error {
	nicknamePart := ctx.QueryParam("nickname")
	if nicknamePart == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var limit uint
	if _, err := fmt.Sscan(ctx.QueryParam("limit"), &limit); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var err error
	var users *[]models.User

	cookie, cookieErr := ctx.Cookie("session_id")
	if cookieErr != nil {
		users, err = userHandler.UseCase.SearchUsers(nil, nicknamePart, limit)
	} else {
		sid := cookie.Value
		uid, sessionErr := userHandler.SessionUseCase.GetByID(sid)
		if sessionErr != nil {
			users, err = userHandler.UseCase.SearchUsers(nil, nicknamePart, limit)
		} else {
			users, err = userHandler.UseCase.SearchUsers(&uid, nicknamePart, limit)
		}
	}

	resp, err := serializer.JSON().Marshal(&users)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
