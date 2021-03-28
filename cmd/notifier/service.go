package main

import (
	"context"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"sport4all/app/handlers/http"
	"sport4all/app/handlers/ws"
	psqlRepos "sport4all/app/repositories/psql"
	redisRepos "sport4all/app/repositories/redis"
	useCases "sport4all/app/usecases/impl"
	"sport4all/pkg/common"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"sync"

	"sport4all/app/handlers/queue"
	"sport4all/app/workers"
	"sport4all/pkg/logger"
	"sport4all/pkg/webSocketPool/gorillaWs"
)

type Service struct {
	receiver    queue.Receiver
	eventWorker workers.EventWorker
	api         ws.Api
}

func CreateService(configFilePath string) *Service {
	settings := InitSettings(configFilePath)
	logger.InitLogger(settings.LogFile, settings.LogLevel)

	rabbitMQOpts := queue.RabbitMQReceiverOpts{
		ConnAddress:       settings.RabbitMQConnAddress,
		QueueId:           settings.RabbitMQEventQueueId,
		MessageBufferSize: 10,
	}

	wsPool := gorillaWs.CreateWebSocketPool()
	logger.Infof("%v", wsPool)

	postgresClient, err := gorm.Open(settings.PsqlName, settings.PsqlData)
	if err != nil {
		logger.Fatal(err)
	}
	defer common.Close(postgresClient.Close)

	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(settings.RedisProtocol, settings.RedisAddress)
			if err != nil {
				logger.Error(err.Error())
			}
			return conn, err
		},
	}

	usrRepo := psqlRepos.CreateUserRepository(postgresClient)
	sessionRepo := redisRepos.CreateSessionRepository(redisPool, settings.RedisExpiresKeySec)

	sesUseCase := useCases.CreateSessionUseCase(sessionRepo, usrRepo)

	receiver, err := queue.NewRabbitMQReceiver(rabbitMQOpts)
	if err != nil {
		logger.Fatal(err)
	}

	mw := http.CreateMiddlewareMini(sesUseCase)

	return &Service{
		receiver:    receiver,
		api:         ws.CreateApi(wsPool, mw),
		eventWorker: workers.CreateEventWorker(wsPool),
		// saver
	}
}

func (s *Service) Run(ctx context.Context) {
	go s.receiver.Run(ctx)
	logger.Info("create receiver worker")

	apiWG := &sync.WaitGroup{}
	apiWG.Add(1)
	go s.api.Run(apiWG)
	logger.Info("create api worker")

	messageWorkerWG := &sync.WaitGroup{}
	messageWorkerWG.Add(1)
	go s.handleMessages(ctx, messageWorkerWG)
	logger.Info("create message worker")

	// add saver

	apiWG.Wait()
	messageWorkerWG.Wait()
	logger.Info("stop service")
}

func (s *Service) handleMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		message, err := s.receiver.TakeMessage(ctx)
		if errors.Is(err, queue.ErrMessageQueueIsClosed) || errors.Is(err, context.Canceled) {
			logger.Info(err)
			return
		}
		if err != nil {
			logger.Error("can't take message: ", err)
			continue
		}

		err = s.eventWorker.ProcessMessage(message)
		if err != nil {
			logger.Error(err)
			continue
		}
	}
}
