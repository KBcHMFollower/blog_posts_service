package repositories_transfer

import "github.com/google/uuid"

type PostConditionTarget string
type PostUpdateTarget string

const (
	PostIdCondition     PostConditionTarget = "id"
	PostUserIdCondition PostConditionTarget = "user_id"
)

const (
	PostTitleUpdateTarget         PostUpdateTarget = "title"
	PostTxtContentUpdateTarget    PostUpdateTarget = "text_content"
	PostImagesContentUpdateTarget PostUpdateTarget = "images_content" //todo: лайки должны быть инкрементом
)

type PostUpdateFieldInfo struct {
	Name  string
	Value string
}

type GetPostsInfo struct {
	Condition map[PostConditionTarget]any
	Size      uint64
	Page      uint64
}

type UpdatePostInfo struct {
	Condition  map[PostConditionTarget]any
	UpdateData map[PostUpdateTarget]any
}

type DeletePostsInfo struct {
	Condition map[PostConditionTarget]any
}

type GetPostInfo struct {
	Condition map[PostConditionTarget]any
}

type GetPostsCountInfo struct {
	Condition map[PostConditionTarget]any
}

type CreatePostInfo struct {
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
}
