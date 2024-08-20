package rabbitmq_client

import (
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/amqp_client"
	"github.com/streadway/amqp"
)

type SendersStore struct {
	ch         *amqp.Channel
	sendersMap map[string]amqp_client.AmqpSender
}

func NewSendersStore(ch *amqp.Channel) *SendersStore {
	sendersMap := map[string]amqp_client.AmqpSender{
		"usersDeleted":   &UserPostsDeletedSender{ch: ch},
		"userCompensate": &UserPostsDeletedSender{ch: ch},
	}

	return &SendersStore{
		ch:         ch,
		sendersMap: sendersMap,
	}
}

func (ss *SendersStore) GetSender(senderName string) (amqp_client.AmqpSender, error) {
	sender, ok := ss.sendersMap[senderName]
	if !ok {
		return nil, fmt.Errorf("sender not found for sender %s", senderName)
	}

	return sender, nil
}

type UserPostsDeletedSender struct {
	ch *amqp.Channel
}

func (upd *UserPostsDeletedSender) Send(message []byte) error {
	if err := upd.ch.Publish(
		DeleteUserExchange,
		UserPostsDeletedQueue,
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

type UserCompensateSender struct {
	ch *amqp.Channel
}

func (uc *UserCompensateSender) Send(message []byte) error {
	if err := uc.ch.Publish(
		DeleteUserExchange,
		UserCompensateQueue,
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
