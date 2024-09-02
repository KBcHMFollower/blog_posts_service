package services_dep

import (
	"context"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type CommentsGetter interface {
	Count(ctx context.Context, postId uuid.UUID) (uint, error)
	Comment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error)
	Comments(ctx context.Context, postId uuid.UUID, size uint64, page uint64) ([]*models.Comment, error)
}

type CommentsCreator interface {
	Create(ctx context.Context, createData repositories_transfer.CreateCommentInfo) (uuid.UUID, error)
}

type CommentsUpdater interface {
	Update(ctx context.Context, updateData repositories_transfer.UpdateCommentInfo) error
}

type CommentsDeleter interface {
	Delete(ctx context.Context, commentId uuid.UUID) error
}
