package queue

import (
	"context"
	"github.com/streadway/amqp"
	"sport4all/app/models"
	"sport4all/pkg/logger"
	"sport4all/pkg/serializer"
)

type RabbitMQReceiverOpts struct {
	ConnAddress       string
	QueueId           string
	MessageBufferSize int
}

type RabbitMQReceiver struct {
	opts         RabbitMQReceiverOpts
	queueConnect *amqp.Connection
	queueChannel *amqp.Channel
	messageBuff  chan *models.Message
}

func NewRabbitMQReceiver(opts RabbitMQReceiverOpts) (Receiver, error) {
	conn, err := amqp.Dial(opts.ConnAddress)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &RabbitMQReceiver{
		opts:         opts,
		queueConnect: conn,
		queueChannel: channel,
		messageBuff:  make(chan *models.Message, opts.MessageBufferSize),
	}, nil
}

func (receiver *RabbitMQReceiver) Run(ctx context.Context) {
	defer receiver.queueConnect.Close()
	defer receiver.queueChannel.Close()

	queue, err := receiver.queueChannel.QueueDeclare(
		receiver.opts.QueueId, // name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		logger.Error(err)
	}

	msgs, err := receiver.queueChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		logger.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var messages []models.Message
			serializer.JSON().Unmarshal(d.Body, &messages)
			for id, _ := range messages {
				receiver.messageBuff <- &messages[id]
			}
			logger.Info("Received a message")
		}
	}()

	logger.Infof("Waiting for messages.")

	<-forever
}

func (receiver *RabbitMQReceiver) TakeMessage(ctx context.Context) (*models.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case mess, ok := <-receiver.messageBuff:
		if !ok {
			return nil, ErrMessageQueueIsClosed
		}
		return mess, nil
	}
}
