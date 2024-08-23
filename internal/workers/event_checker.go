package workers

import (
	"context"
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
	workers_dep "github.com/KBcHMFollower/blog_posts_service/internal/workers/interfaces/dep"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type EventStore interface {
	workers_dep.EventGetter
	workers_dep.EventUpdater
}

type EventChecker struct {
	amqpClient amqpclient.AmqpClient //TODO: РАЗБИТЬ НА СЕНДЕР И КОНСЬЮМЕР
	eventRep   EventStore
	logger     *slog.Logger
}

func NewEventChecker(amqpClient amqpclient.AmqpClient, eventRep EventStore, logger *slog.Logger) *EventChecker {
	return &EventChecker{
		amqpClient: amqpClient,
		eventRep:   eventRep,
		logger:     logger,
	}
}

func (as EventChecker) Run(ctx context.Context) error {
	log := as.logger.
		With("op", "EventChecker.Run")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				events, err := as.eventRep.GetEvents(ctx, "status", "waiting", 50)
				if err != nil {
					log.Error("can`t get events from  db: ", err)
					time.Sleep(5 * time.Second)
					continue
				}

				for _, event := range events {

					err = as.amqpClient.Publish(event.EventType, []byte(event.Payload))
					if err != nil {
						log.Error("can`t publish event: ", err)
						time.Sleep(5 * time.Second)
						continue
					}
				}

				eventIds := make([]uuid.UUID, 0, len(events))
				for _, event := range events {
					eventIds = append(eventIds, event.EventId)
				}

				err = as.eventRep.SetSentStatusesInEvents(ctx, eventIds)
				if err != nil {
					log.Error("can`t set status in events: ", err)
					time.Sleep(5 * time.Second)
					continue
				}

				time.Sleep(5 * time.Second)
			}

		}
	}()

	return nil
}
