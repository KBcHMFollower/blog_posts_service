package services

import (
	"context"
	"fmt"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	"github.com/KBcHMFollower/blog_posts_service/internal/repository"
	services_dep "github.com/KBcHMFollower/blog_posts_service/internal/services/interfaces/dep"
	"log/slog"
)

type RequestsStore interface {
	services_dep.RequestsCreator
	services_dep.RequestsGetter
}

type RequestsService struct {
	reqRepository repository.RequestsStore
	log           *slog.Logger
}

func NewRequestsService(reqRepository repository.RequestsStore, log *slog.Logger) *RequestsService {
	return &RequestsService{
		reqRepository: reqRepository,
		log:           log,
	}
}

func (rs *RequestsService) CheckExists(ctx context.Context, checkInfo services_transfer.RequestsCheckExistsInfo) (bool, error) {
	op := "PostService.GetUserPosts"

	log := rs.log.With(
		slog.String("op", op),
	)

	res, err := rs.reqRepository.Get(ctx, checkInfo.Key)
	if err != nil {
		log.Error(err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if res == nil {
		return false, nil
	}

	return true, nil
}
