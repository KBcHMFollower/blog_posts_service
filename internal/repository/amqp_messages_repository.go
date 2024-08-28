package repository

import (
	"context"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	rep_utils "github.com/KBcHMFollower/blog_posts_service/internal/repository/lib"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	messagesStatusCol    = "status"
	messagesEventIdCol   = "event_id"
	messagesAllCol       = "*"
	messagesEventTypeCol = "event_type"
	messagesPayloadCol   = "payload"
)

const (
	SentStatus = "sent"
)

type EventFilter struct {
}

type EventRepository struct {
	db       database.DBWrapper
	qBuilder squirrel.StatementBuilderType //TODO: СЛИШКОМ ПРЯМАЯ ЗАВИСИМОСТЬ
}

func NewEventRepository(dbDriver database.DBWrapper) *EventRepository {
	return &EventRepository{
		db:       dbDriver,
		qBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *EventRepository) GetEvents(ctx context.Context, filterTarget string, filterValue interface{}, limit uint64) ([]*models.EventInfo, error) {
	op := "UserRepository.getSubInfo"

	query := r.qBuilder.
		Select(messagesAllCol).
		From(database.AmqpMessagesTable).
		Where(squirrel.Eq{filterTarget: filterValue}).
		Limit(limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.GenerateSqlErr(err, op)
	}

	eventInfos := make([]*models.EventInfo, 0)
	err = r.db.SelectContext(ctx, &eventInfos, sql, args...)
	if err != nil {
		return nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return eventInfos, nil
}

func (r *EventRepository) SetSentStatusesInEvents(ctx context.Context, eventsId []uuid.UUID) error {
	op := "UserRepository.getSubInfo"

	query := r.qBuilder.
		Update(database.AmqpMessagesTable).
		Where(squirrel.Eq{messagesEventIdCol: eventsId}).
		Set(messagesStatusCol, SentStatus)

	sql, args, err := query.ToSql()
	if err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	return nil
}

func (r *EventRepository) GetEventById(ctx context.Context, eventId uuid.UUID) (*models.EventInfo, error) {
	op := "UserRepository.getEventById"

	query := r.qBuilder.
		Select(messagesAllCol).
		From(database.AmqpMessagesTable).
		Where(squirrel.Eq{messagesEventIdCol: eventId})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.GenerateSqlErr(err, op)
	}

	var eventInfo models.EventInfo
	if err := r.db.GetContext(ctx, &eventInfo, sql, args...); err != nil {
		return nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return &eventInfo, nil
}

func (r *EventRepository) Create(ctx context.Context, info repositories_transfer.CreateEventInfo, tx database.Transaction) error {
	op := "UserRepository.create"
	executor := rep_utils.GetExecutor(r.db, tx)

	query := r.qBuilder.
		Insert(database.AmqpMessagesTable).
		SetMap(map[string]interface{}{
			messagesEventTypeCol: info.EventId,
			messagesEventIdCol:   info.EventId,
			messagesPayloadCol:   info.Payload,
		})

	toSql, args, err := query.ToSql()
	if err != nil {
		return rep_utils.GenerateSqlErr(err, op)
	}

	_, err = executor.ExecContext(ctx, toSql, args...)
	if err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	return nil
}
