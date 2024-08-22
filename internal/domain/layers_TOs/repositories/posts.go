package repositories_transfer

import "github.com/google/uuid"

type PostUpdateFieldInfo struct {
	Name  string
	Value string
}

type GetPostByUserIdInfo struct {
	UserId uuid.UUID
	Size   uint32
	Page   uint32
}

type UpdatePostInfo struct {
	Id         uuid.UUID
	UpdateData []*CommentUpdateFieldInfo
}

type CreatePostInfo struct {
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
}
