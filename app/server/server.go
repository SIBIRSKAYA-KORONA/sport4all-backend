package server

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"

	httpHandlers "github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/handlers/http"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	psqlRepos "github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories/psql"
	redisRepos "github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories/redis"
	useCases "github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases/impl"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/common"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
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

	// postgresClient.DropTableIfExists(&models.User{})
	postgresClient.AutoMigrate(&models.User{})

	usrRepo := psqlRepos.CreateUserRepository(postgresClient)

	/* USE CASES */
	sesUseCase := useCases.CreateSessionUseCase(sessionRepo, usrRepo)
	usrUseCase := useCases.CreateUserUseCase(sessionRepo, usrRepo)

	/* HANDLERS */
	mw := httpHandlers.CreateMiddleware(sesUseCase)
	router := echo.New()
	router.Use(mw.ProcessPanic)
	router.Use(mw.LogRequest)
	//router.Use(mw.CORS)
	router.Use(mw.Sanitize)
	rootGroup := router.Group(server.settings.BaseURL)

	httpHandlers.CreateSessionHandler(server.settings.SessionsURL, rootGroup, sesUseCase, mw)
	httpHandlers.CreateUserHandler(server.settings.SettingsURL, server.settings.ProfileURL, rootGroup, usrUseCase, mw)

	logger.Error("start server on address: ", server.settings.ServerAddress,
		", log file: ", server.settings.LogFile, ", log level: ", server.settings.LogLevel)

	if err = router.Start(server.settings.ServerAddress); err != nil {
		logger.Fatal(err)
	}
}
