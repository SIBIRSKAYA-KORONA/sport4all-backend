package queue

import (
	"context"
	"errors"

	"sport4all/app/models"
)

var ErrMessageQueueIsClosed = errors.New("message queue is closed")

type Receiver interface {
	Run(ctx context.Context)
	TakeMessage(ctx context.Context) (*models.Message, error)
}
