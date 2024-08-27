package services

import (
	"context"
	"fmt"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
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

// TODO: ДОБАВИТЬ ИНТЕРЦЕПТОРЫ С ЭТИМ МЕТОДОМ
func (rs *RequestsService) CheckAndCreate(ctx context.Context, checkInfo services_transfer.RequestsCheckExistsInfo) (bool, error) {
	op := "PostService.GetUserPosts"

	log := rs.log.With(
		slog.String("op", op),
	)

	res, err := rs.reqRepository.Get(ctx, checkInfo.Key, nil)
	if err != nil {
		log.Error(err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if res != nil {
		return true, nil
	}

	_, _, err = rs.reqRepository.Create(ctx, repositories_transfer.CreateRequestInfo{
		Key: checkInfo.Key,
	}, nil)
	if err != nil {
		log.Error(err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return false, nil
}
