package services_dep

import (
	"context"
	"database/sql"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type EventSetter interface {
	SetSentStatusesInEvents(ctx context.Context, eventsId []uuid.UUID) error
}

type EventGetter interface {
	GetEvents(ctx context.Context, filterTarget string, filterValue interface{}, limit uint64) ([]*models.EventInfo, error)
	GetEventById(ctx context.Context, eventId uuid.UUID) (*models.EventInfo, error)
}

type EventCreator interface {
	Create(ctx context.Context, info repositories_transfer.CreateEventInfo, tx *sql.Tx) error
}
