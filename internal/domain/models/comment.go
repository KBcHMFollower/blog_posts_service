package models

import (
	commentsv1 "github.com/KBcHMFollower/test_plate_blog_service/api/protos/gen/comments"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Comment struct {
	Id        uuid.UUID
	PostId    uuid.UUID
	UserId    uuid.UUID
	Content   string
	Likes     uint32
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
		Id:        c.Id.String(),
		UserId:    c.UserId.String(),
		PostId:    c.PostId.String(),
		Content:   c.Content,
		Likes:     int32(c.Likes),
		CreatedAt: timestamppb.New(c.CreatedAt),
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
