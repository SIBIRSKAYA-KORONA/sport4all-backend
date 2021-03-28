package workers

import (
	"sport4all/app/models"
	"sport4all/pkg/logger"
)

type EventWorker interface {
	ProcessMessage(message *models.Message) error
}

func CreateEventWorker() EventWorker {
	return &EventWorkerImpl{}
}

type EventWorkerImpl struct {
	// websockets pool ?
}

func (worker *EventWorkerImpl) ProcessMessage (message *models.Message) error {
	logger.Debugf("Got message: %v", *message)
	logger.Debug("Do some work with message")
	return nil
}


