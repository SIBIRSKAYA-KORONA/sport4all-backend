package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"sport4all/pkg/sanitize"
	"sport4all/pkg/serializer"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"

	"sport4all/app/models"
	useCases "sport4all/app/usecases"
	"sport4all/pkg/common"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type Middleware interface {
	LogRequest(echo.HandlerFunc) echo.HandlerFunc
	ProcessPanic(echo.HandlerFunc) echo.HandlerFunc
	Sanitize(echo.HandlerFunc) echo.HandlerFunc
	DebugMiddle(echo.HandlerFunc) echo.HandlerFunc
	CORS(echo.HandlerFunc) echo.HandlerFunc
	CheckAuth(echo.HandlerFunc) echo.HandlerFunc
	CheckTeamPermission(role models.Role) echo.MiddlewareFunc
	CheckTournamentPermission(role models.TournamentRole) echo.MiddlewareFunc
	CheckTournamentPermissionByMeeting(role models.TournamentRole) echo.MiddlewareFunc
	CheckMeetingStatus(status models.EventStatus) echo.MiddlewareFunc
	CheckTeamInMeeting(echo.HandlerFunc) echo.HandlerFunc
	CheckPlayerInTeam() echo.MiddlewareFunc
	NotificationMiddleware(models.MessageType) echo.MiddlewareFunc
}

type MiddlewareImpl struct {
	sessionUseCase    useCases.SessionUseCase
	teamUseCase       useCases.TeamUseCase
	tournamentUseCase useCases.TournamentUseCase
	mettingUseCase    useCases.MeetingUseCase
	messageUseCase    useCases.MessageUseCase
	origins           map[string]struct{}

	attachURL string

	channel *amqp.Channel
	queue   amqp.Queue
}

func CreateMiddleware(sessionUseCase useCases.SessionUseCase,
	teamUseCase useCases.TeamUseCase,
	tournamentUseCase useCases.TournamentUseCase,
	meetingUseCase useCases.MeetingUseCase,
	messageUseCase useCases.MessageUseCase,
	origins map[string]struct{},
	attachURL string,
	channel *amqp.Channel,
	queue amqp.Queue) Middleware {
	return &MiddlewareImpl{
		sessionUseCase:    sessionUseCase,
		teamUseCase:       teamUseCase,
		tournamentUseCase: tournamentUseCase,
		mettingUseCase:    meetingUseCase,
		messageUseCase:    messageUseCase,
		origins:           origins,
		attachURL:         attachURL,
		channel:           channel,
		queue:             queue,
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
		if (ctx.Request().Method != echo.PUT && ctx.Request().Method != echo.POST) || ctx.Path() == mw.attachURL {
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
		if _, exist := mw.origins[origin]; len(origin) != 0 && !exist {
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

func (mw *MiddlewareImpl) DebugMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		dump, err := httputil.DumpRequest(ctx.Request(), true)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		logger.Debugf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump)
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
		return mw.CheckAuth(func(ctx echo.Context) error {
			var teamID uint
			if _, err := fmt.Sscan(ctx.Param("tid"), &teamID); err != nil {
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
			}

			return ctx.String(errors.ResolveErrorToCode(errors.ErrNoPermission), errors.ErrNoPermission.Error())
		})
	}
}

func (mw *MiddlewareImpl) CheckTournamentPermission(role models.TournamentRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return mw.CheckAuth(func(ctx echo.Context) error {
			var tournamentID uint
			if _, err := fmt.Sscan(ctx.Param("tournamentId"), &tournamentID); err != nil {
				return ctx.NoContent(http.StatusBadRequest)
			}
			userID := ctx.Get("uid").(uint)

			if _, err := mw.tournamentUseCase.GetByID(tournamentID); err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}
			ctx.Set("tournamentId", tournamentID)

			ok, err := mw.tournamentUseCase.CheckUserForTournamentRole(tournamentID, userID, role)
			if err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}

			if ok {
				return next(ctx)
			}

			return ctx.String(errors.ResolveErrorToCode(errors.ErrNoPermission), errors.ErrNoPermission.Error())
		})
	}
}

func (mw *MiddlewareImpl) CheckTournamentPermissionByMeeting(role models.TournamentRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return mw.CheckAuth(func(ctx echo.Context) error {
			var meetingID uint
			if _, err := fmt.Sscan(ctx.Param("mid"), &meetingID); err != nil {
				return ctx.NoContent(http.StatusBadRequest)
			}
			userID := ctx.Get("uid").(uint)

			meeting, err := mw.mettingUseCase.GetByID(meetingID)
			if err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}

			ctx.Set("tournamentId", meeting.TournamentId)
			ctx.Set("meetingId", meetingID)

			ok, err := mw.tournamentUseCase.CheckUserForTournamentRole(meeting.TournamentId, userID, role)
			if err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}

			if ok {
				return next(ctx)
			}

			return ctx.String(errors.ResolveErrorToCode(errors.ErrNoPermission), errors.ErrNoPermission.Error())
		})
	}
}

func (mw *MiddlewareImpl) CheckMeetingStatus(status models.EventStatus) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return mw.CheckAuth(func(ctx echo.Context) error {
			var meetingId uint
			if _, err := fmt.Sscan(ctx.Param("mid"), &meetingId); err != nil {
				return ctx.NoContent(http.StatusBadRequest)
			}

			meeting, err := mw.mettingUseCase.GetByID(meetingId)
			if err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}

			userId := ctx.Get("uid").(uint)

			ok, err := mw.tournamentUseCase.CheckUserForTournamentRole(meeting.TournamentId, userId, models.TournamentOrganizer)
			if err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}

			if !ok {
				return ctx.String(errors.ResolveErrorToCode(errors.ErrNoPermission), errors.ErrNoPermission.Error())
			}

			if meeting.Status != status {
				return ctx.String(errors.ResolveErrorToCode(errors.ErrMeetingStatusNotAcceptable), errors.ErrMeetingStatusNotAcceptable.Error())
			}

			ctx.Set("meetingId", meeting.ID)
			ctx.Set("tournamentId", meeting.TournamentId)

			return next(ctx)
		})
	}
}

func (mw *MiddlewareImpl) CheckTeamInMeeting(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		meetingId := ctx.Get("meetingId").(uint)

		var teamId uint
		if _, err := fmt.Sscan(ctx.Param("tid"), &teamId); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		result, err := mw.mettingUseCase.IsTeamInMeeting(meetingId, teamId)
		if err != nil {
			logger.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}

		if !result {
			return ctx.String(errors.ResolveErrorToCode(errors.ErrMeetingNotFound), errors.ErrMeetingNotFound.Error())
		}

		ctx.Set("teamId", teamId)
		return next(ctx)
	}
}

func (mw *MiddlewareImpl) CheckPlayerInTeam() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return mw.CheckTeamInMeeting(func(ctx echo.Context) error {
			var playerId uint
			if _, err := fmt.Sscan(ctx.Param("uid"), &playerId); err != nil {
				return ctx.NoContent(http.StatusBadRequest)
			}

			teamId := ctx.Get("teamId").(uint)

			result, err := mw.teamUseCase.CheckUserForRole(teamId, playerId, models.Player)
			if err != nil {
				logger.Error(err)
				return ctx.String(errors.ResolveErrorToCode(err), err.Error())
			}

			if !result {
				return ctx.String(errors.ResolveErrorToCode(errors.ErrNoPermission), errors.ErrNoPermission.Error())
			}

			ctx.Set("playerId", playerId)
			return next(ctx)
		})
	}
}

func (mw *MiddlewareImpl) NotificationMiddleware(messageType models.MessageType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)
			status := ctx.Response().Status
			if err != nil || status != http.StatusOK {
				logger.Error("error:", err, " status:", status)
				return err
			}

			messages := mw.fillMessageByType(ctx, messageType)
			//logger.Info(len(*messages))

			if err := mw.messageUseCase.Create(messages); err != nil {
				logger.Error(err)
			}

			encoded, err := serializer.JSON().Marshal(&messages)
			if err != nil {
				logger.Error(err)
				return ctx.NoContent(http.StatusInternalServerError)
			}

			err = mw.channel.Publish(
				"",            // exchange
				mw.queue.Name, // routing key
				false,         // mandatory
				false,         // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        encoded,
				})
			if err != nil {
				logger.Error(err)
			}

			return next(ctx)
		}
	}
}

func (mw *MiddlewareImpl) fillMessageByType(ctx echo.Context, messageType models.MessageType) *[]models.Message {
	var messages []models.Message

	switch messageType {

	// -----------------------------------------------------------
	case models.MeetingStatusChanged:
		// получаем встречу
		teams, err := mw.tournamentUseCase.GetAllTeams(ctx.Get("tournamentId").(uint))
		if err != nil {
			logger.Error(err)
			return nil
		}

		messagesByUser := make(map[uint]bool)
		status := ctx.Get("status").(uint)

		var messageType models.MessageType
		if status == 2 {
			messageType = models.MeetingStarted
		} else if status == 3 {
			messageType = models.MeetingFinished
		}

		meetingId := ctx.Get("meetingId").(uint)
		// собираем всех игроков, которым будем отправлять уведомление
		for teamID, _ := range *teams {
			for _, player := range (*teams)[teamID].Players {
				_, alreadySent := messagesByUser[player.ID]
				if !alreadySent {
					message := models.Message{
						MessageType: messageType,
						TargetUid:   player.ID,
						MeetingId:   meetingId,
						CreateAt:    time.Now().Unix(),
						IsRead:      false,
					}

					messages = append(messages, message)
					messagesByUser[player.ID] = true
				}
			}

			teamOwnerId := (*teams)[teamID].OwnerId
			_, alreadySent := messagesByUser[teamOwnerId]
			if !alreadySent {
				ownerMessage := models.Message{
					MessageType: messageType,
					TargetUid:   teamOwnerId,
					MeetingId:   meetingId,
					CreateAt:    time.Now().Unix(),
					IsRead:      false,
				}
				messages = append(messages, ownerMessage)
				messagesByUser[teamOwnerId] = true
			}
		}

	case models.AddedToTeam:
		message := models.Message{
			MessageType: models.AddedToTeam,
			TargetUid:   ctx.Get("member").(uint),
			SourceUid:   ctx.Get("uid").(uint),
			MeetingId:   0,
			CreateAt:    time.Now().Unix(),
			IsRead:      false,
		}
		messages = append(messages, message)
	}

	return &messages
}
