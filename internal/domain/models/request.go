package models

import (
	"github.com/google/uuid"
)

type Request struct {
	Id             uuid.UUID
	IdempotencyKey uuid.UUID
	Status         string
}

func (r *Request) GetPointersArray() []interface{} {
	return []interface{}{
		&r.Id,
		&r.IdempotencyKey,
		&r.Status,
	}
}
