package rabbitmqclient

import (
	"fmt"
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
	"github.com/streadway/amqp"
)

type SendersStore struct {
	ch         *amqp.Channel
	sendersMap map[string]amqpclient.AmqpSender
}

func NewSendersStore(ch *amqp.Channel) *SendersStore {
	sendersMap := map[string]amqpclient.AmqpSender{
		"postsDeleted": &PostsDeletedSender{ch: ch},
	}

	return &SendersStore{
		ch:         ch,
		sendersMap: sendersMap,
	}
}

func (ss *SendersStore) GetSender(senderName string) (amqpclient.AmqpSender, error) {
	sender, ok := ss.sendersMap[senderName]
	if !ok {
		return nil, fmt.Errorf("sender not found for sender %s", senderName)
	}

	return sender, nil
}

type PostsDeletedSender struct {
	ch *amqp.Channel
}

func (s *PostsDeletedSender) Send(message []byte) error {
	if err := s.ch.Publish(
		DeleteUserExchange,
		UserDeletedQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	); err != nil {
		return fmt.Errorf("failed to send message: %s", err)
	}

	return nil
}
