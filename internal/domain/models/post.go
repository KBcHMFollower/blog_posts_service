package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
	Likes         int32
	CreatedAt     time.Time
}

func CreatePost(userId uuid.UUID, title string, textContent string, imageContent *string) *Post {
	return &Post{
		Id:            uuid.New(),
		UserId:        userId,
		Title:         title,
		TextContent:   textContent,
		ImagesContent: imageContent,
	}
}
