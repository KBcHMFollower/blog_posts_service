package repositories_transfer

import "github.com/google/uuid"

type CommentUpdateFieldInfo struct {
	Name  string
	Value string
}

type UpdateCommentInfo struct {
	Id         uuid.UUID
	UpdateData []*CommentUpdateFieldInfo
}

type CreateCommentInfo struct {
	PostId  uuid.UUID
	UserId  uuid.UUID
	Content string
}
