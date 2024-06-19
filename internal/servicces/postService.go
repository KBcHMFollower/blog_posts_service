package postService

import "github.com/KBcHMFollower/test_plate_user_service/internal/repository"

type PostService struct {
	postRepository repository.PostRepository
}

func New(postRepository repository.PostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}
