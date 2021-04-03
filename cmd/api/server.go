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
	"github.com/streadway/amqp"
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
	/*  Redis */
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

	/* PostgreSQL */
	postgresClient, err := gorm.Open(server.settings.PsqlName, server.settings.PsqlData)
	if err != nil {
		logger.Fatal(err)
	}
	defer common.Close(postgresClient.Close)

	s3session, err := session.NewSession(&aws.Config{Region: aws.String(server.settings.S3Region)})
	if err != nil {
		logger.Fatal(err)
	}

	postgresClient.AutoMigrate(&models.User{}, &models.Team{}, &models.Tournament{}, &models.Meeting{},
		&models.Stats{}, &models.Attach{}, &models.Skill{}, &models.SkillApprove{})

	/* RabbitMQ */
	conn, err := amqp.Dial(server.settings.RabbitMQConnAddress)
	if err != nil {
		logger.Fatal(err)
	}
	defer common.Close(conn.Close)

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal(err)
	}
	defer common.Close(ch.Close)

	queue, err := ch.QueueDeclare(
		server.settings.RabbitMQEventQueueId, // name
		false,                                // durable
		false,                                // delete when unused
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)
	if err != nil {
		logger.Fatal(err)
	}

	/* REPOS */
	sessionRepo := redisRepos.CreateSessionRepository(redisPool, server.settings.RedisExpiresKeySec)

	userRepo := psqlRepos.CreateUserRepository(postgresClient)
	teamRepo := psqlRepos.CreateTeamRepository(postgresClient)
	tournamentRepo := psqlRepos.CreateTournamentRepository(postgresClient)
	meetingRepo := psqlRepos.CreateMeetingRepository(postgresClient)
	skillRepo := psqlRepos.CreateSkillRepository(postgresClient)
	attachRepo := amazonS3Repos.CreateAttachRepository(postgresClient, s3session, server.settings.S3Bucket)

	/* USE CASES */
	sesUseCase := useCases.CreateSessionUseCase(sessionRepo, userRepo)
	usrUseCase := useCases.CreateUserUseCase(sessionRepo, userRepo)
	teamUseCase := useCases.CreateTeamUseCase(teamRepo, userRepo)
	tournamentUseCase := useCases.CreateTournamentUseCase(userRepo, tournamentRepo, teamRepo, meetingRepo)
	meetingUseCase := useCases.CreateMeetingUseCase(meetingRepo, tournamentRepo)
	skillUseCase := useCases.CreateSkillUseCase(skillRepo, userRepo)
	attachUseCase := useCases.CreateAttachUseCase(attachRepo)

	/* HANDLERS */
	origins := make(map[string]struct{})
	for _, key := range server.settings.Origins {
		origins[key] = struct{}{}
	}

	mw := httpHandlers.CreateMiddleware(sesUseCase, teamUseCase, tournamentUseCase, meetingUseCase, origins,
		server.settings.BaseURL+server.settings.AttachURL, ch, queue)

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
	httpHandlers.CreateSkillHandler(server.settings.SkillsURL, rootGroup, skillUseCase, mw)
	httpHandlers.CreateAttachHandler(server.settings.AttachURL, rootGroup, attachUseCase, mw)

	logger.Info("start server on address: ", server.settings.ServerAddress,
		", log file: ", server.settings.LogFile, ", log level: ", server.settings.LogLevel)

	if err = router.Start(server.settings.ServerAddress); err != nil {
		logger.Fatal(err)
	}
}
