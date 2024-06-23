package repository

import (
	"context"

	"github.com/KBcHMFollower/test_plate_blog_service/internal/domain/models"
	"github.com/google/uuid"
)

const (
	POSTS_TABLE    = "posts"
	COMMENTS_TABLE = "comments"
	ID_FIELD       = "id"
	USER_ID_FIELD  = "user_id"
)

type UpdateItem struct {
	Name  string
	Value string
}

type UpdateData struct {
	Id         uuid.UUID
	UpdateData []*UpdateItem
}

type CreatePostData struct {
	User_id       uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
}

type CreateCommentData struct {
	PostId  uuid.UUID
	UserId  uuid.UUID
	Content string
}

type PostStore interface {
	CreatePost(ctx context.Context, createData CreatePostData) (uuid.UUID, *models.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	GetPostsByUserId(ctx context.Context, user_id uuid.UUID, size uint64, page uint64) ([]*models.Post, uint, error)
	DeletePost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	UpdatePost(ctx context.Context, updateData UpdateData) (*models.Post, error)
}

type CommentStore interface {
	CreateComment(ctx context.Context, createData CreateCommentData) (uuid.UUID, *models.Comment, error)
	GetComment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error)
	GetPostComments(ctx context.Context, postId uuid.UUID, size uint64, page uint64) ([]*models.Comment, uint, error)
	DeleteComment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error)
	UpdateComment(ctx context.Context, updateData UpdateData) (*models.Comment, error)
}
