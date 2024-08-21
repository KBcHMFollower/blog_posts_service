package services

import (
	"context"
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/repository"
	services_dep "github.com/KBcHMFollower/blog_posts_service/internal/services/interfaces/dep"
	"github.com/google/uuid"
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

func (rs *RequestsService) CheckExists(ctx context.Context, key uuid.UUID, payload string) (bool, error) {
	op := "PostService.GetUserPosts"

	log := rs.log.With(
		slog.String("op", op),
	)

	res, err := rs.reqRepository.Get(ctx, key)
	if err != nil {
		log.Error(err.Error())
		return true, fmt.Errorf("%s: %w", op, err)
	}

	if res != nil {
		return true, nil
	}

	_, _, err = rs.reqRepository.Create(ctx, key, payload)
	if err != nil {
		log.Error(err.Error())
		return true, fmt.Errorf("%s: %w", op, err)
	}

	return false, nil
}
