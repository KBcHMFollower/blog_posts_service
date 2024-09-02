package services_dep

import (
	"context"
	"database/sql"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type PostCreator interface {
	Create(ctx context.Context, createData repositories_transfer.CreatePostInfo) (uuid.UUID, error)
}

type PostGetter interface {
	Post(ctx context.Context, id uuid.UUID) (*models.Post, error)
	GetPostsByUserId(ctx context.Context, getInfo repositories_transfer.GetPostsInfo) ([]*models.Post, error)
	GetUserPostsCount(ctx context.Context, userId uuid.UUID) (uint, error)
}

type PostDeleter interface {
	DeletePost(ctx context.Context, id uuid.UUID) error
	DeleteUserPosts(ctx context.Context, userId uuid.UUID, tx *sql.Tx) error
}

type PostUpdater interface {
	UpdatePost(ctx context.Context, updateData repositories_transfer.UpdatePostInfo) (*models.Post, error)
}
