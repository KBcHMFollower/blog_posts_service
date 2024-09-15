package services_transfer

import (
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
)

type GetCommentInfo struct {
	CommId uuid.UUID
}

type GetPostCommentsInfo struct {
	PostId uuid.UUID
	Size   uint64
	Page   uint64
}

type DeleteCommentInfo struct {
	CommId uuid.UUID
}

type UpdateCommentInfo struct {
	CommId       uuid.UUID
	UpdateFields map[repositories_transfer.CommentUpdateTarget]any
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
	Likes   uint64
}

type GetPostCommentsResult struct {
	TotalCount uint64
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
		Likes:   model.Likes,
	}
}

func ConvertCommentsArrayFromModels(models []*models.Comment) []CommentResult {
	var results []CommentResult = make([]CommentResult, 0, len(models))

	for _, model := range models {
		results = append(results, ConvertCommentFromModel(model))
	}

	return results
}

func (c *CommentResult) ToProto() *commentsv1.Comment {
	return &commentsv1.Comment{
		Id:      c.CommId.String(),
		PostId:  c.PostId.String(),
		UserId:  c.UserId.String(),
		Content: c.Content,
		Likes:   c.Likes,
	}
}

func CommentsArrayProto(c []CommentResult) []*commentsv1.Comment {
	var res = make([]*commentsv1.Comment, 0, len(c))

	for _, comment := range c {
		res = append(res, comment.ToProto())
	}

	return res
}
