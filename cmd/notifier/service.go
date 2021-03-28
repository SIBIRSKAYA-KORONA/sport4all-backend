package main

import (
	"context"
	"errors"
	"sync"

	"sport4all/app/receiver"
	"sport4all/app/workers"
	"sport4all/pkg/logger"
)

type Service struct {
	receiver  receiver.Receiver
	eventWorker workers.EventWorker
}

func CreateService(configFilePath string) *Service {
	settings := InitSettings(configFilePath)
	logger.InitLogger(settings.LogFile, settings.LogLevel)


	rabbitMQOpts := receiver.RabbitMQReceiverOpts{
		ConnAddress: settings.RabbitMQConnAddress,
		QueueId: settings.RabbitMQEventQueueId,
		MessageBufferSize: 10,
	}

	receiver, err := receiver.NewRabbitMQReceiver(rabbitMQOpts)
	if err != nil {
		logger.Fatal(err)
	}

	return &Service{
		receiver: receiver,
		eventWorker: workers.CreateEventWorker(),
	}
}

func (s *Service) Run(ctx context.Context) {

	go s.receiver.Run(ctx)
	logger.Info("create receiver worker")


	messageWorkerWG := &sync.WaitGroup{}
	messageWorkerWG.Add(1)
	go s.handleMessages(ctx, messageWorkerWG)


	messageWorkerWG.Wait()
	logger.Info("stop service")
}

func (s *Service) handleMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		message, err := s.receiver.TakeMessage(ctx)
		if errors.Is(err, receiver.ErrMessageQueueIsClosed) || errors.Is(err, context.Canceled) {
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


