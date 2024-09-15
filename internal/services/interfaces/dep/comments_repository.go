package services_dep

import (
	"context"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type CommentsGetter interface {
	Count(ctx context.Context, info repositories_transfer.GetCommentsCountInfo) (uint64, error)
	Comment(ctx context.Context, info repositories_transfer.GetCommentInfo) (*models.Comment, error)
	Comments(ctx context.Context, info repositories_transfer.GetCommentsInfo) ([]*models.Comment, error)
}

type CommentsCreator interface {
	Create(ctx context.Context, createData repositories_transfer.CreateCommentInfo) (uuid.UUID, error)
}

type CommentsUpdater interface {
	Update(ctx context.Context, updateData repositories_transfer.UpdateCommentInfo) error
}

type CommentsDeleter interface {
	Delete(ctx context.Context, info repositories_transfer.DeleteCommentInfo) error
}
