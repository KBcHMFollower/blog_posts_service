package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrBadRequest          = errors.New("bad request")
	ErrForbidden           = errors.New("forbidden")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrConflict            = errors.New("conflict")
	ErrInternalServerError = errors.New("internal server error")
)

func AddOpInErr(err error, op string) error {
	return fmt.Errorf("%s : %w", op, err)
}
