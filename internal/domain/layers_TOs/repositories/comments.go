package repositories_transfer

import "github.com/google/uuid"

type UpdateItem struct {
	Name  string
	Value string
}

type UpdateData struct {
	Id         uuid.UUID
	UpdateData []*UpdateItem
}

type CreatePostData struct {
	User_id       uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
}

type CreateCommentData struct {
	PostId  uuid.UUID
	UserId  uuid.UUID
	Content string
}
