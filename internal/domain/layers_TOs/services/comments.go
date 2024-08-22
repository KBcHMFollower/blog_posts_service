package services_transfer

import (
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type CommUpdateFieldInfo struct {
	Name  string
	Value string
}

type GetCommentInfo struct {
	CommId uuid.UUID
}

type GetPostCommentsInfo struct {
	PostId uuid.UUID
	Size   int32
	Page   int32
}

type DeleteCommentInfo struct {
	CommId uuid.UUID
}

type UpdateCommentInfo struct {
	CommId       uuid.UUID
	UpdateFields []CommUpdateFieldInfo
}

type CreateCommentInfo struct {
	UserId  uuid.UUID
	PostId  uuid.UUID
	Content string
} //TODO: тут вроде еще же картинковый контент может быть

type CommentResult struct {
	CommId  uuid.UUID
	PostId  uuid.UUID
	UserId  uuid.UUID
	Content string
	Likes   int32
}

type GetPostCommentsResult struct {
	TotalCount int32
	Comments   []CommentResult
}

type GetCommentResult struct {
	Comment CommentResult
}

type UpdateCommentResult struct {
	CommId  uuid.UUID
	Comment CommentResult
}

type CreateCommentResult struct {
	CommId  uuid.UUID
	Comment CommentResult
}

func ConvertCommentFromModel(model *models.Comment) CommentResult {
	return CommentResult{
		CommId:  model.Id,
		PostId:  model.PostId,
		UserId:  model.UserId,
		Content: model.Content,
		Likes:   int32(model.Likes),
	}
}

func ConvertCommentsArrayFromModels(models []*models.Comment) []CommentResult {
	var results []CommentResult = make([]CommentResult, 0, len(models))

	for _, model := range models {
		results = append(results, ConvertCommentFromModel(model))
	}

	return results
}
