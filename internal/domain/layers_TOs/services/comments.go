package services_transfer

import (
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (c *CommentResult) ToProto() *commentsv1.Comment {
	return &commentsv1.Comment{
		Id:        c.CommId.String(),
		PostId:    c.PostId.String(),
		UserId:    c.UserId.String(),
		Content:   c.Content,
		Likes:     int32(c.Likes),
		CreatedAt: &timestamppb.Timestamp{},
	} //TODO: дата зоздания елиента ебать не должна
}

func CommentsArrayProto(c []CommentResult) []*commentsv1.Comment {
	var res = make([]*commentsv1.Comment, 0, len(c))

	for _, comment := range c {
		res = append(res, comment.ToProto())
	}

	return res
}
