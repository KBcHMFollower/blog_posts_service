package services_dep

import (
	"context"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type RequestsCreator interface {
	Create(ctx context.Context, info repositories_transfer.CreateRequestInfo) (uuid.UUID, *models.Request, error)
}

type RequestsGetter interface {
	Get(ctx context.Context, key uuid.UUID) (*models.Request, error)
}
