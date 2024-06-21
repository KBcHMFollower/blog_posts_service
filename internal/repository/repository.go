package repository

import (
	"context"

	"github.com/KBcHMFollower/test_plate_blog_service/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateItem struct {
	Name  string
	Value string
}

type UpdatePostData struct {
	Id         uuid.UUID
	UpdateData []*UpdateItem
}

type CreatePostData struct {
	User_id       uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
}

type IPostRepository interface {
	CreatePost(ctx context.Context, createData CreatePostData) (uuid.UUID, *models.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	GetPostsByUserId(ctx context.Context, user_id uuid.UUID, size uint64, page uint64) ([]*models.Post, uint, error)
	DeletePost(ctx context.Context, id uuid.UUID) (*models.Post, error)
	UpdatePost(ctx context.Context, updateData UpdatePostData) (*models.Post, error)
}
