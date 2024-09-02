package repositories_transfer

import "github.com/google/uuid"

type CommentConditionTarget string
type CommentUpdateTarget string

const (
	CommentPostIdConditionTarget CommentConditionTarget = "post_id"
	CommentUserIdConditionTarget CommentConditionTarget = "user_id"
	CommentIdConditionTarget     CommentConditionTarget = "id"
)

type UpdateCommentInfo struct {
	Condition  map[CommentConditionTarget]interface{}
	UpdateData map[CommentUpdateTarget]interface{}
}

type GetCommentInfo struct {
	Condition map[CommentConditionTarget]interface{}
}

type GetCommentsCountInfo struct {
	Condition map[CommentConditionTarget]interface{}
}

type GetCommentsInfo struct {
	Condition map[CommentConditionTarget]interface{}
	Page      uint64
	Size      uint64
}

type DeleteCommentInfo struct {
	Condition map[CommentConditionTarget]interface{}
}

type CreateCommentInfo struct {
	PostId  uuid.UUID
	UserId  uuid.UUID
	Content string
}
