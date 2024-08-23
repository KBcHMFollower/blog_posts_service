package services_dep

import (
	"context"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type PostCreator interface {
	CreatePost(ctx context.Context, createData repositories_transfer.CreatePostInfo) (uuid.UUID, *models.Post, error)
}

type PostGetter interface {
	GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	GetPostsByUserId(ctx context.Context, getInfo repositories_transfer.GetPostByUserIdInfo) ([]*models.Post, uint, error)
}

type PostDeleter interface {
	DeletePost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	DeleteUserPosts(ctx context.Context, userId uuid.UUID) error
}

type PostUpdater interface {
	UpdatePost(ctx context.Context, updateData repositories_transfer.UpdatePostInfo) (*models.Post, error)
}
