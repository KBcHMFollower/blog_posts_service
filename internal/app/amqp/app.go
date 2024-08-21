package amqpapp

import (
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
)

type AmqpApp struct {
	client   amqpclient.AmqpClient
	handlers map[string]amqpclient.AmqpHandlerFunc
}

func NewAmqpApp(client amqpclient.AmqpClient) *AmqpApp {
	return &AmqpApp{
		client:   client,
		handlers: make(map[string]amqpclient.AmqpHandlerFunc),
	}
}

func (app *AmqpApp) RegisterHandler(name string, handler amqpclient.AmqpHandlerFunc) {
	app.handlers[name] = handler
}

func (app *AmqpApp) Start() error {
	for name, handler := range app.handlers {
		err := app.client.Consume(name, handler)
		if err != nil {
			return err
		}
	}

	return nil
} //TODO : STOP
