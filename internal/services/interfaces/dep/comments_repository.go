package services_dep

import (
	"context"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type CommentsGetter interface {
	GetPostCommentsCount(ctx context.Context, postId uuid.UUID) (uint, error)
	GetComment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error)
	GetPostComments(ctx context.Context, postId uuid.UUID, size uint64, page uint64) ([]*models.Comment, error)
}

type CommentsCreator interface {
	CreateComment(ctx context.Context, createData repositories_transfer.CreateCommentInfo) (uuid.UUID, error)
}

type CommentsUpdater interface {
	UpdateComment(ctx context.Context, updateData repositories_transfer.UpdateCommentInfo) error
}

type CommentsDeleter interface {
	DeleteComment(ctx context.Context, commentId uuid.UUID) error
}
