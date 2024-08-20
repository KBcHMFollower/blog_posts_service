package amqp_client

type AmqpSender interface {
	Send(message []byte) error
}

type AmqpConsumer interface {
	StartConsume() error
}

type AmqpSenderFactory interface {
	GetSender(eventType string) (AmqpSender, error)
}

type AmqpClient interface {
	GetSendersProvider() AmqpSenderFactory
	GetConsumer() AmqpConsumer
}
