package models

import (
	"fmt"
	ssov1 "github.com/KBcHMFollower/test_plate_blog_service/api/protos/gen/posts"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Post struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent string
	Likes         int32
	CreatedAt     time.Time
}

func CreatePost(userId uuid.UUID, title string, textContent string, imageContent string) *Post {
	return &Post{
		Id:            uuid.New(),
		UserId:        userId,
		Title:         title,
		TextContent:   textContent,
		ImagesContent: imageContent,
	}
}

func (p *Post) ConvertToProto() *ssov1.Post {
	return &ssov1.Post{
		Id:            p.Id.String(),
		UserId:        p.UserId.String(),
		Title:         p.Title,
		TextContent:   p.TextContent,
		ImagesContent: p.ImagesContent,
		Likes:         p.Likes,
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
}

func ConvertFromProto(protoPost *ssov1.Post) (*Post, error) {

	postUUID, err := uuid.Parse(protoPost.GetId())
	if err != nil {
		return nil, fmt.Errorf("can`t parse post_id in uuid: %v", err)
	}

	userUUID, err := uuid.Parse(protoPost.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("can`t parse user_id in uuid: %v", err)
	}

	return &Post{
		Id:            postUUID,
		UserId:        userUUID,
		Title:         protoPost.GetTitle(),
		TextContent:   protoPost.GetTextContent(),
		ImagesContent: protoPost.ImagesContent,
		Likes:         protoPost.GetLikes(),
		CreatedAt:     protoPost.GetCreatedAt().AsTime(),
	}, nil
}

func (p *Post) GetPointersArray() []interface{} {
	return []interface{}{
		&p.Id,
		&p.UserId,
		&p.Title,
		&p.TextContent,
		&p.ImagesContent,
		&p.Likes,
		&p.CreatedAt,
	}
}
