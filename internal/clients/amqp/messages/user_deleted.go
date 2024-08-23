package messages

import "github.com/google/uuid"

type UserDeletedMessage struct {
	EventId uuid.UUID
	UserId  uuid.UUID
}
