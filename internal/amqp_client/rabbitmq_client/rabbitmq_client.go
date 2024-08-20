package rabbitmq_client

import (
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/amqp_client"
	"github.com/streadway/amqp"
)

const (
	DeleteUserExchange    = "direct-user-actions"
	UserDeletedQueue      = "user-deleted"
	UserPostsDeletedQueue = "user-posts-deleted"
	UserCompensateQueue   = "user-compensate"
)

type RabbitMQClient struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	sendersStore amqp_client.AmqpSenderFactory
	consumer     amqp_client.AmqpConsumer
}

func NewRabbitMQClient(addr string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	sendersFactory := NewSendersStore(ch)
	consumer := NewConsumer(ch)

	return &RabbitMQClient{ch: ch, conn: conn, sendersStore: sendersFactory, consumer: consumer}, nil
}

func (rc *RabbitMQClient) Close() error {
	if err := rc.ch.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %s", err)
	}
	if err := rc.conn.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection: %s", err)
	}

	return nil
}

func (rc *RabbitMQClient) GetSendersProvider() amqp_client.AmqpSenderFactory {
	return rc.sendersStore
}

func (rc *RabbitMQClient) GetConsumer() amqp_client.AmqpConsumer {
	return rc.consumer
}
