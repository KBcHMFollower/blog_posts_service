package rabbitmq

import "github.com/streadway/amqp"

const (
	POST_EXCHANGE_NAME     = "exchange.posts"
	POST_DELETE_QUEUE_NAME = "queue.posts.delete"
)

type Connection struct {
	channel    *amqp.Channel
	connection *amqp.Connection
}

func New() (*Connection, error) {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Connection{ch, conn}, nil
}

func DeclareExchangeForPosts(conn *Connection) error {
	err := conn.channel.ExchangeDeclare(
		POST_EXCHANGE_NAME, // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func DeclareAndBindDeletePostQueue(conn *Connection) error {
	q, err := conn.channel.QueueDeclare(
		POST_DELETE_QUEUE_NAME, // name (пустое имя позволяет RabbitMQ создать временную очередь)
		false,                  // durable
		false,                  // delete when unused
		true,                   // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return err
	}

	err = conn.channel.QueueBind(
		q.Name,             // queue name
		q.Name,             // routing key
		POST_EXCHANGE_NAME, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func PublishDeletePostMessage(connection *Connection, userId string) error {
	err := connection.channel.Publish(
		POST_EXCHANGE_NAME,
		POST_DELETE_QUEUE_NAME,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(userId),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
