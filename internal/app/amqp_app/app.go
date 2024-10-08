package amqpapp

import (
	"fmt"
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
	"github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp/rabbitmqclient"
	"github.com/KBcHMFollower/blog_posts_service/internal/config"
	ctxerrors "github.com/KBcHMFollower/blog_posts_service/internal/domain/errors"
	"github.com/KBcHMFollower/blog_posts_service/internal/logger"
)

type AmqpApp struct {
	Client   amqpclient.AmqpClient
	handlers map[string]amqpclient.AmqpHandlerFunc
}

func NewAmqpApp(rabbitmqConnectInfo config.RabbitMq, log logger.Logger) (*AmqpApp, error) {
	rabbitMqApp, err := rabbitmqclient.NewRabbitMQClient(rabbitmqConnectInfo.Addr, log)
	if err != nil {
		return nil, ctxerrors.Wrap("can`t to connect to rabbitmq", err)
	}

	return &AmqpApp{
		Client:   rabbitMqApp,
		handlers: make(map[string]amqpclient.AmqpHandlerFunc),
	}, nil
}

func (app *AmqpApp) RegisterHandler(name string, handler amqpclient.AmqpHandlerFunc) {
	app.handlers[name] = handler
}

func (app *AmqpApp) Start() error {
	for name, handler := range app.handlers {
		err := app.Client.Consume(name, handler)
		if err != nil {
			return ctxerrors.Wrap(fmt.Sprintf("error in subscribe to query `%s`", name), err)
		}
	}

	return nil
}

func (app *AmqpApp) Stop() error {
	if err := app.Client.Stop(); err != nil {
		return ctxerrors.Wrap("can`t to stop rabbitmq", err)
	}
	return nil
}
