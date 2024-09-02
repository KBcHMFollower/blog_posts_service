package handlers_utils

import (
	"fmt"
	ctxerrors "github.com/KBcHMFollower/blog_posts_service/internal/domain/errors"
)

func ReturnValidationError(err error) error {
	return fmt.Errorf("%w:\n%s", ctxerrors.ErrBadRequest, err.Error())
}
