package messages

import (
	"github.com/google/uuid"
	"time"
)

type PostsDeleteFeedbackMessage struct {
	UserId uuid.UUID `json:"user_id"`
	Ok     bool      `json:"ok"`
}

type UserDeletedMessage struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FName       string    `json:"fname"`
	LName       string    `json:"lname"`
	Avatar      string    `json:"avatar"`
	AvatarMin   string    `json:"avatar_min"`
	IsDeleted   bool      `json:"is_deleted"`
	PassHash    []byte    `json:"pass_hash"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}
