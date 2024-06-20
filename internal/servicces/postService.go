package postService

import (
	"log/slog"

	"github.com/KBcHMFollower/test_plate_user_service/internal/repository"
)

type PostService struct {
	postRepository repository.PostRepository
	log            *slog.Logger
}

func New(postRepository repository.PostRepository, log *slog.Logger) *PostService {
	return &PostService{
		postRepository: postRepository,
		log:            log,
	}
}
