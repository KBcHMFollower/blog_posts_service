package services_transfer

import (
	postsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/posts"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateUserFieldInfo struct {
	Name  string
	Value string
}

type GetUserPostsInfo struct {
	UserId uuid.UUID
	Size   int32
	Page   int32
}

type GetPostInfo struct {
	PostId uuid.UUID
}

type DeletePostInfo struct {
	PostId uuid.UUID
}

type CreatePostInfo struct {
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
}

type UpdatePostInfo struct {
	PostId uuid.UUID
	Fields []UpdateUserFieldInfo
}

type DeleteUserPostInfo struct {
	UserId uuid.UUID
}

type PostResult struct {
	PostId        uuid.UUID
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
	Likes         int32
}

type GetUserPostsResult struct {
	Posts      []PostResult
	TotalCount int32
}

type GetPostResult struct {
	Post PostResult
}

type CreatePostResult struct {
	PostId uuid.UUID
	Post   PostResult
}

type UpdatePostResult struct {
	PostId uuid.UUID
	Post   PostResult
}

func ConvertPostFromModel(model *models.Post) PostResult {
	return PostResult{
		PostId:        model.Id,
		UserId:        model.UserId,
		Title:         model.Title,
		TextContent:   model.TextContent,
		ImagesContent: model.ImagesContent,
		Likes:         model.Likes,
	}
}

func ConvertPostsArrayFromModel(posts []*models.Post) []PostResult {
	results := make([]PostResult, 0, len(posts))

	for _, post := range posts {
		results = append(results, ConvertPostFromModel(post))
	}

	return results
}

func (p *PostResult) ToProto() *postsv1.Post {
	return &postsv1.Post{
		Id:            p.PostId.String(),
		UserId:        p.UserId.String(),
		Title:         p.Title,
		TextContent:   p.TextContent,
		ImagesContent: p.ImagesContent,
		Likes:         p.Likes,
		CreatedAt:     nil,
	} //TODO: CreatedAt не нужно
}

func ConvertPostArrayToProto(posts []PostResult) []*postsv1.Post {
	results := make([]*postsv1.Post, 0, len(posts))
	for _, post := range posts {
		results = append(results, post.ToProto())
	}
	return results
}
