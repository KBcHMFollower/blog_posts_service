package services_transfer

import "github.com/google/uuid"

type UpdateUserFieldInfo struct {
	Name  string
	Value string
}

type GetUserPostsInfo struct {
	UserId uuid.UUID
	Size   int32
	Page   int32
}

type GetPostInfo struct {
	PostId uuid.UUID
}

type DeletePostInfo struct {
	PostId uuid.UUID
}

type CreatePostInfo struct {
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
}

type UpdatePostInfo struct {
	PostId uuid.UUID
	Fields []UpdateUserFieldInfo
}

type DeleteUserPostInfo struct {
	UserId uuid.UUID
}

type PostResult struct {
	PostId        uuid.UUID
	UserId        uuid.UUID
	Title         string
	TextContent   string
	ImagesContent *string
	Likes         int32
}

type GetUserPostsResult struct {
	Posts      []PostResult
	TotalCount int32
}

type GetPostResult struct {
	Post PostResult
}

type CreatePostResult struct {
	PostId uuid.UUID
	Post   PostResult
}

type UpdatePostResult struct {
	PostId uuid.UUID
	Post   PostResult
}
