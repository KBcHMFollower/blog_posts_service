package models

import (
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Id        uuid.UUID
	PostId    uuid.UUID
	UserId    uuid.UUID
	Content   string
	Likes     uint64
	CreatedAt time.Time
}

func CreateComment(postId uuid.UUID, userId uuid.UUID, content string) *Comment {
	return &Comment{
		Id:      uuid.New(),
		PostId:  postId,
		UserId:  userId,
		Content: content,
		Likes:   0,
	}
}

func (c *Comment) ConvertToProto() *commentsv1.Comment {
	return &commentsv1.Comment{
		Id:      c.Id.String(),
		UserId:  c.UserId.String(),
		PostId:  c.PostId.String(),
		Content: c.Content,
		Likes:   c.Likes,
	}
}

func (c *Comment) GetPointerArray() []interface{} {
	return []interface{}{
		&c.Id,
		&c.PostId,
		&c.UserId,
		&c.Content,
		&c.Likes,
		&c.CreatedAt,
	}
}
