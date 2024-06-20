package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id           uuid.UUID
	User_id      uuid.UUID
	Title        string
	TextContent  string
	ImageContent string
	Created_at   time.Time
}

func CreatePost(user_id uuid.UUID, title string, textContent string, imageContent string) *Post {
	return &Post{
		Id:           uuid.New(),
		User_id:      user_id,
		Title:        title,
		TextContent:  textContent,
		ImageContent: imageContent,
		Created_at:   time.Now(),
	}
}
