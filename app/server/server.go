package server

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"

	handlers "github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/handlers/http"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/middleware"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories/postgreSQL"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases/impl"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/config"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type Server struct {
	ip             string
	port           uint
	configObserver config.Observer
}

func CreateServer() *Server {
	configObserver_ := config.CreateConfigObserver()
	return &Server{
		ip:             configObserver_.GetServerIP(),
		port:           configObserver_.GetServerPort(),
		configObserver: configObserver_,
	}
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.ip, server.port)
}

func (server *Server) Run() {
	router := echo.New()
	rootGroup := router.Group(server.configObserver.GetBaseURL())
	postgresClient, err := gorm.Open(server.configObserver.GetDBMS(), server.configObserver.GetDBConnection())
	if err != nil {
		logger.Fatal(err)
	}
	defer postgresClient.Close()
	logger.Info(server.configObserver.GetDBMS() + " initialized ✓")

	postgresClient.AutoMigrate(&models.User{})

	usrRepo := postgreSQL.CreateUserRepository(postgresClient)
	logger.Info("entities repositories initialized ✓")

	usrUseCase := impl.CreateUserUseCase(usrRepo)
	logger.Info("entities usecases initialized ✓")

	mw := middleware.CreateMiddleware()
	router.Use(mw.ProcessPanic)
	router.Use(mw.LogRequest)
	//router.Use(mw.CORS)
	logger.Info("middlewares initialized ✓")

	handlers.CreateUserHandler(rootGroup, server.configObserver.GetSettingsURL(), usrUseCase, mw)
	logger.Info("entities handlers initialized ✓")

	err = router.Start(server.GetAddr())
	if err != nil {
		logger.Fatal(err)
	}
}
