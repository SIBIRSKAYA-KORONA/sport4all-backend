package main

import (
	httpHandlers "sport4all/app/handlers/http"
	"sport4all/app/models"
	amazonS3Repos "sport4all/app/repositories/amazon_s3"
	psqlRepos "sport4all/app/repositories/psql"
	redisRepos "sport4all/app/repositories/redis"
	useCases "sport4all/app/usecases/impl"
	"sport4all/pkg/common"
	"sport4all/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
)

type Server struct {
	settings Settings
}

func CreateServer(configFilePath string) *Server {
	settings := InitSettings(configFilePath)
	logger.InitLogger(settings.LogFile, settings.LogLevel)
	return &Server{settings: settings}
}

func (server *Server) Run() {
	/* REPOS */
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(server.settings.RedisProtocol, server.settings.RedisAddress)
			if err != nil {
				logger.Error(err.Error())
			}
			return conn, err
		},
	}
	defer common.Close(redisPool.Close)

	sessionRepo := redisRepos.CreateSessionRepository(redisPool, server.settings.RedisExpiresKeySec)

	postgresClient, err := gorm.Open(server.settings.PsqlName, server.settings.PsqlData)
	if err != nil {
		logger.Fatal(err)
	}
	defer common.Close(postgresClient.Close)

	s3session, err := session.NewSession(&aws.Config{Region: aws.String(server.settings.S3Region)})
	if err != nil {
		logger.Fatal(err)
	}

	// postgresClient.DropTableIfExists(&models.User{}, &models.Team{}, &models.Tournament{}, &models.Meeting{},
	//&models.Stats{}, &models.Attach{})
	postgresClient.AutoMigrate(&models.User{}, &models.Team{}, &models.Tournament{}, &models.Meeting{},
		&models.Stats{}, &models.Attach{})

	usrRepo := psqlRepos.CreateUserRepository(postgresClient)
	teamRepo := psqlRepos.CreateTeamRepository(postgresClient)
	tournamentRepo := psqlRepos.CreateTournamentRepository(postgresClient)
	meetingRepo := psqlRepos.CreateMeetingRepository(postgresClient)
	attachRepo := amazonS3Repos.CreateAttachRepository(postgresClient, s3session, server.settings.S3Bucket)

	/* USE CASES */
	sesUseCase := useCases.CreateSessionUseCase(sessionRepo, usrRepo)
	usrUseCase := useCases.CreateUserUseCase(sessionRepo, usrRepo)
	teamUseCase := useCases.CreateTeamUseCase(teamRepo, usrRepo)
	tournamentUseCase := useCases.CreateTournamentUseCase(usrRepo, tournamentRepo, teamRepo, meetingRepo)
	meetingUseCase := useCases.CreateMeetingUseCase(meetingRepo, tournamentRepo)
	attachUseCase := useCases.CreateAttachUseCase(attachRepo)

	/* HANDLERS */
	origins := make(map[string]struct{})
	for _, key := range server.settings.Origins {
		origins[key] = struct{}{}
	}

	mw := httpHandlers.CreateMiddleware(sesUseCase, teamUseCase, tournamentUseCase, meetingUseCase, origins, server.settings.AttachURL)
	router := echo.New()
	router.Use(mw.ProcessPanic)
	router.Use(mw.LogRequest)
	router.Use(mw.CORS)
	router.Use(mw.Sanitize)
	rootGroup := router.Group(server.settings.BaseURL)

	httpHandlers.CreateSessionHandler(server.settings.SessionsURL, rootGroup, sesUseCase, mw)
	httpHandlers.CreateUserHandler(server.settings.SettingsURL, server.settings.ProfileURL, rootGroup, usrUseCase, mw)
	httpHandlers.CreateTeamHandler(server.settings.TeamsURL, rootGroup, teamUseCase, mw)
	httpHandlers.CreateTournamentHandler(server.settings.TournamentsURL, rootGroup, tournamentUseCase, mw)
	httpHandlers.CreateMeetingsHandler(server.settings.MeetingsURL, rootGroup, meetingUseCase, mw)
	httpHandlers.CreateAttachHandler(server.settings.AttachURL, rootGroup, attachUseCase, mw)

	logger.Info("start server on address: ", server.settings.ServerAddress,
		", log file: ", server.settings.LogFile, ", log level: ", server.settings.LogLevel)

	if err = router.Start(server.settings.ServerAddress); err != nil {
		logger.Fatal(err)
	}
}
