package rabbitmqclient

import (
	"context"
	"errors"
	"fmt"
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
	ctxerrors "github.com/KBcHMFollower/blog_posts_service/internal/domain/errors"
	"github.com/KBcHMFollower/blog_posts_service/internal/logger"
	"github.com/streadway/amqp"
)

const (
	DeleteUserExchange    = "direct-user-actions"
	UserDeletedQueue      = amqpclient.UserDeletedEventKey
	UserPostsDeletedQueue = amqpclient.PostsDeletedEventKey
)

const (
	queueLogKey   = "queue"
	messageLogKey = "message"
)

type RabbitMQClient struct {
	pubConn        *amqp.Connection
	pubCh          *amqp.Channel
	consConn       *amqp.Connection
	consCh         *amqp.Channel
	sendersFactory amqpclient.AmqpSenderFactory
	log            logger.Logger
	ctx            context.Context
}

func NewRabbitMQClient(addr string, log logger.Logger) (amqpclient.AmqpClient, error) {
	ctx := context.Background()
	ctx = logger.UpdateLoggerCtx(ctx, "worker-name", "rabbitmq")

	pConn, err := amqp.Dial(addr)
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("failed to connect to RabbitMQ", err))
	}
	cConn, err := amqp.Dial(addr)
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("failed to connect to RabbitMQ", err))
	}

	pCh, err := pConn.Channel()
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("failed to open a channel", err))
	}
	cCh, err := cConn.Channel()
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("failed to open a channel", err))
	}

	if err := DeclareExchanges(pCh); err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("failed to declare exchanges", err))
	}
	if err := DeclareQueues(pCh); err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("failed to declare queues", err))
	}

	sendersFactory := NewSendersStore(pCh)

	return &RabbitMQClient{ctx: context.Background(), pubConn: cConn, consConn: cConn, pubCh: pCh, consCh: cCh, sendersFactory: sendersFactory, log: log}, nil
}

func (rc *RabbitMQClient) Close() error {
	if err := rc.pubCh.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %s", err)
	}
	if err := rc.pubConn.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection: %s", err)
	}
	if err := rc.consCh.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %s", err)
	}
	if err := rc.consConn.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection: %s", err)
	}

	return nil
}

func (rc *RabbitMQClient) Publish(ctx context.Context, eventType string, body []byte) error {
	sender, err := rc.sendersFactory.GetSender(eventType)
	if err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap(fmt.Sprintf("failed to get sender for event %s", eventType), err))
	}

	if err := sender.Send(body); err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap(fmt.Sprintf("failed to send event %s", eventType), err))
	}

	return nil
}

func (rc *RabbitMQClient) Consume(queueName string, handler amqpclient.AmqpHandlerFunc) error {
	del, err := rc.consCh.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return ctxerrors.WrapCtx(rc.ctx, ctxerrors.Wrap(fmt.Sprintf("failed to register a consumer for queue %s", queueName), err))
	}

	rc.log.InfoContext(rc.ctx, "start consuming messages", queueLogKey, queueName)

	go func() {

		ctx, cancel := context.WithCancel(rc.ctx)
		ctx = logger.UpdateLoggerCtx(ctx, queueLogKey, queueName)

		for d := range del {
			select {
			case <-rc.ctx.Done():
				cancel()
				return
			default:

				rc.log.InfoContext(ctx, "received a message", messageLogKey, string(d.Body))
				if err := handler(ctx, d.Body); err != nil {
					rc.log.ErrorContext(ctxerrors.ErrorCtx(ctx, err), "failed to handle a message from queue", logger.ErrKey, err.Error())
					if err := d.Nack(false, false); err != nil {
						rc.log.ErrorContext(ctxerrors.ErrorCtx(ctx, err), "failed to nack", logger.ErrKey, err.Error())
					}
					continue
				}
				rc.log.InfoContext(ctx, "finished consuming a message from queue")
				if err := d.Ack(false); err != nil {
					rc.log.ErrorContext(ctxerrors.ErrorCtx(ctx, err), "failed to ack", logger.ErrKey, err)
				}
			}

		}
		cancel()
	}()

	return nil
}

func (rc *RabbitMQClient) Stop() error { //TODO: МОЖНО В STOP CTX ПРОКИДЫВАТЬ ДЛЯ ЛОГГЕРА
	var resErr error = nil

	if err := rc.pubCh.Close(); err != nil {
		resErr = errors.Join(resErr, fmt.Errorf("failed to close RabbitMQ channel: %s", err))
	}
	if err := rc.pubConn.Close(); err != nil {
		resErr = errors.Join(resErr, fmt.Errorf("failed to close RabbitMQ connection: %s", err))
	}
	if err := rc.consCh.Close(); err != nil {
		resErr = errors.Join(resErr, fmt.Errorf("failed to close RabbitMQ channel: %s", err))
	}
	if err := rc.consConn.Close(); err != nil {
		resErr = errors.Join(resErr, fmt.Errorf("failed to close RabbitMQ connection: %s", err))
	}

	rc.ctx.Done()

	return resErr
}

func DeclareExchanges(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		DeleteUserExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare DeleteUser exchange: %s", err)
	}

	return nil
}

func DeclareQueues(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		UserDeletedQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare UserDeleted queue: %s", err)
	}

	if err = ch.QueueBind(
		q.Name,
		UserDeletedQueue,
		DeleteUserExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind UserDeleted queue: %s", err)
	}

	q, err = ch.QueueDeclare(
		UserPostsDeletedQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare UserDeleted queue: %s", err)
	}

	if err = ch.QueueBind(
		q.Name,
		UserPostsDeletedQueue,
		DeleteUserExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind UserDeleted queue: %s", err)
	}

	return nil
}
