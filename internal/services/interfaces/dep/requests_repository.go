package services_dep

import (
	"context"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type RequestsCreator interface {
	Create(ctx context.Context, key uuid.UUID, payload string) (uuid.UUID, *models.Request, error)
}

type RequestsGetter interface {
	Get(ctx context.Context, key uuid.UUID) (*models.Request, error)
}
