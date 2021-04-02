package workers

import (
	"sport4all/app/models"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"
	"sport4all/pkg/webSocketPool"
)

type EventWorker interface {
	ProcessMessage(message *models.Message) error
}

func CreateEventWorker(wsPool webSocketPool.WebSocketPool) EventWorker {
	return &EventWorkerImpl{wsPool: wsPool}
}

type EventWorkerImpl struct {
	wsPool webSocketPool.WebSocketPool
}

func (worker *EventWorkerImpl) ProcessMessage(message *models.Message) error {
	logger.Debugf("Got message: %v", *message)
	logger.Debug("Do some work with message")

	resp, err := serializer.JSON().Marshal(message)
	if err != nil {
		logger.Error(err)
		return err
	}

	worker.wsPool.Send(message.TargetUid, resp)
	return nil
}
