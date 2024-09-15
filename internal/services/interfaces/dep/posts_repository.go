package services_dep

import (
	"context"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type PostCreator interface {
	Create(ctx context.Context, createData repositories_transfer.CreatePostInfo) (uuid.UUID, error)
}

type PostGetter interface {
	Post(ctx context.Context, info repositories_transfer.GetPostInfo) (*models.Post, error)
	Posts(ctx context.Context, getInfo repositories_transfer.GetPostsInfo) ([]*models.Post, error)
	Count(ctx context.Context, info repositories_transfer.GetPostsCountInfo) (uint64, error)
}

type PostDeleter interface {
	Delete(ctx context.Context, info repositories_transfer.DeletePostsInfo, tx database.Transaction) error
}

type PostUpdater interface {
	Update(ctx context.Context, updateData repositories_transfer.UpdatePostInfo) error
}
