package rabbitmq_client

import (
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/services/post_service"
	"github.com/KBcHMFollower/blog_posts_service/internal/services/requests_service"
	"github.com/streadway/amqp"
)

type Consumer struct {
	reqService   requests_service.RequestsService
	postsService post_service.PostService
	ch           *amqp.Channel
}

func NewConsumer(ch *amqp.Channel, reqService requests_service.RequestsService, postsService post_service.PostService) *Consumer {
	return &Consumer{
		ch:           ch,
		reqService:   reqService,
		postsService: postsService,
	}
}

func (c *Consumer) StartConsume() error {
	if err := c.consumeUserDeletedEvent(); err != nil {
		return fmt.Errorf("error consuming user deleted event: %s", err)
	}

	return nil
}

func (c *Consumer) consumeUserDeletedEvent() error {
	deliveries, err := c.defaultConsumeQueue(UserDeletedQueue)
	if err != nil {
		return fmt.Errorf("error getting deliveries: %s", err)
	}

	for d := range deliveries {
		
	}
}

func (c *Consumer) defaultConsumeQueue(queue string) (<-chan amqp.Delivery, error) {
	deliveries, err := c.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error on consume %v", err)
	}

	return deliveries, nil
}
