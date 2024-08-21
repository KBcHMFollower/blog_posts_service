package services_dep

import (
	"context"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type PostCreator interface {
	CreatePost(ctx context.Context, createData repositories_transfer.CreatePostData) (uuid.UUID, *models.Post, error)
}

type PostGetter interface {
	GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	GetPostsByUserId(ctx context.Context, user_id uuid.UUID, size uint64, page uint64) ([]*models.Post, uint, error)
}

type PostDeleter interface {
	DeletePost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	DeleteUserPosts(ctx context.Context, userId uuid.UUID) error
}

type PostUpdater interface {
	UpdatePost(ctx context.Context, updateData repositories_transfer.UpdateData) (*models.Post, error)
}
