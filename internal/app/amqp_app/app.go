package amqpapp

import (
	"fmt"
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
	"github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp/rabbitmqclient"
	"github.com/KBcHMFollower/blog_posts_service/internal/config"
)

type AmqpApp struct {
	Client   amqpclient.AmqpClient
	handlers map[string]amqpclient.AmqpHandlerFunc
}

func New(rabbitmqConnectInfo config.RabbitMq) (*AmqpApp, error) {
	rabbitMqApp, err := rabbitmqclient.NewRabbitMQClient(rabbitmqConnectInfo.Addr)
	if err != nil {
		return nil, fmt.Errorf("new rabbitmq Client error: %v", err)
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
			return err
		}
	}

	return nil
}

func (app *AmqpApp) Stop() error {
	return app.Client.Close()
}
