package services_transfer

import "github.com/google/uuid"

type RequestsCheckExistsInfo struct {
	key     uuid.UUID
	payload string
}
