package services

import (
	"context"
	"encoding/json"
	"fmt"
	amqpclient "github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp"
	"github.com/KBcHMFollower/blog_posts_service/internal/clients/amqp/messages"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	services_dep "github.com/KBcHMFollower/blog_posts_service/internal/services/interfaces/dep"
	"github.com/google/uuid"
	"log/slog"
)

type PostsStore interface {
	services_dep.PostCreator
	services_dep.PostDeleter
	services_dep.PostGetter
	services_dep.PostUpdater
}

type RequestStore interface {
	services_dep.RequestsCreator
	services_dep.RequestsGetter
}

type EventStore interface {
	services_dep.EventCreator
}

type PostService struct {
	postRepository     PostsStore
	requestsRepository RequestStore
	eventsRepository   EventStore
	txCreator          services_dep.TransactionCreator
	log                *slog.Logger
}

func NewPostsService(
	postRepository PostsStore,
	requestsRepository RequestStore,
	eventsRepository EventStore,
	txCreator services_dep.TransactionCreator,
	log *slog.Logger,
) *PostService {
	return &PostService{
		postRepository:     postRepository,
		log:                log,
		eventsRepository:   eventsRepository,
		requestsRepository: requestsRepository,
		txCreator:          txCreator,
	}
}

func (g *PostService) GetUserPosts(ctx context.Context, getInfo *services_transfer.GetUserPostsInfo) (*services_transfer.GetUserPostsResult, error) {
	op := "PostService.GetUserPosts"
	log := g.log.With(
		slog.String("op", op),
	)

	posts, err := g.postRepository.GetPostsByUserId(ctx, repositories_transfer.GetPostsInfo{
		UserId: getInfo.UserId,
		Size:   uint32(getInfo.Size),
		Page:   uint32(getInfo.Page),
	})
	if err != nil {
		log.Error(fmt.Sprintf("failed to get posts for user %d: %v", getInfo.UserId, err))
		return nil, domain.AddOpInErr(err, op)
	}
	total, err := g.postRepository.GetUserPostsCount(ctx, getInfo.UserId)
	if err != nil {
		log.Error(fmt.Sprintf("failed to get posts count for user %d: %v", getInfo.UserId, err))
		return nil, domain.AddOpInErr(err, op)
	}

	return &services_transfer.GetUserPostsResult{
		Posts:      services_transfer.ConvertPostsArrayFromModel(posts),
		TotalCount: int32(total),
	}, nil
}

func (g *PostService) GetPost(ctx context.Context, getInfo *services_transfer.GetPostInfo) (*services_transfer.GetPostResult, error) {
	op := "PostService.GetPost"
	log := g.log.With(
		slog.String("op", op),
	)

	post, err := g.postRepository.Post(ctx, getInfo.PostId)
	if err != nil {
		log.Error(fmt.Sprintf("failed to get post for post %d: %v", getInfo.PostId, err))
		return nil, domain.AddOpInErr(err, op)
	}

	return &services_transfer.GetPostResult{
		Post: services_transfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) DeletePost(ctx context.Context, deleteInfo *services_transfer.DeletePostInfo) error {
	op := "PostService.Delete"

	log := g.log.With(
		slog.String("op", op),
	)

	if err := g.postRepository.DeletePost(ctx, deleteInfo.PostId); err != nil {
		log.Error(fmt.Sprintf("failed to delete post %d: %v", deleteInfo.PostId, err))
		return domain.AddOpInErr(err, op)
	}

	return nil
}

func (g *PostService) CreatePost(ctx context.Context, createInfo *services_transfer.CreatePostInfo) (*services_transfer.CreatePostResult, error) {
	op := "PostService.Create"
	log := g.log.With(
		slog.String("op", op),
	)

	postId, err := g.postRepository.Create(ctx, repositories_transfer.CreatePostInfo{
		UserId:        createInfo.UserId,
		Title:         createInfo.Title,
		TextContent:   createInfo.TextContent,
		ImagesContent: createInfo.ImagesContent,
	})
	if err != nil {
		log.Error(fmt.Sprintf("failed to create post for user %d: %v", createInfo.UserId, err))
		return nil, domain.AddOpInErr(err, op)
	}
	post, err := g.postRepository.Post(ctx, postId)
	if err != nil {
		log.Error(fmt.Sprintf("failed to get post for post %d: %v", postId, err))
		return nil, domain.AddOpInErr(err, op)
	}

	return &services_transfer.CreatePostResult{
		PostId: postId,
		Post:   services_transfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) UpdatePost(ctx context.Context, updateInfo *services_transfer.UpdatePostInfo) (*services_transfer.UpdatePostResult, error) {
	op := "PostService.Create"

	log := g.log.With(
		slog.String("op", op),
	)

	updateItems := make([]*repositories_transfer.CommentUpdateFieldInfo, 0)

	for _, item := range updateInfo.Fields {
		updateItems = append(updateItems, &repositories_transfer.CommentUpdateFieldInfo{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	post, err := g.postRepository.UpdatePost(ctx, repositories_transfer.UpdatePostInfo{
		Id:         updateInfo.PostId,
		UpdateData: updateItems,
	})
	if err != nil {
		log.Error(fmt.Sprintf("failed to update post for post %d: %v", updateInfo.PostId, err))
		return nil, domain.AddOpInErr(err, op)
	}

	return &services_transfer.UpdatePostResult{
		PostId: post.Id,
		Post:   services_transfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) DeleteUserPosts(ctx context.Context, deleteInfo services_transfer.DeleteUserPostInfo) (resErr error) { //TODO: СДЕЛАТЬ DEFER
	op := "PostsService.deleteUserPosts"

	tx, err := g.txCreator.BeginTx(ctx, nil)
	if err != nil {
		g.log.Error("can`t begin transaction :", err)
		return domain.AddOpInErr(g.createPostDeleteFeedbackEvent(
			ctx,
			messages.PostsDeleted{Status: messages.Failed, EventId: deleteInfo.EventId},
			err,
		), op)
	}
	defer func() {

	}()
	//TODO: ОСТАНОВИЛСЯ ЗДЕСЬ
	if err := g.postRepository.DeleteUserPosts(ctx, deleteInfo.UserId, tx); err != nil {
		g.log.Error("can`t delete user from db :", err)
		rbErr := tx.Rollback()
		return domain.AddOpInErr(g.createPostDeleteFeedbackEvent(
			ctx,
			messages.PostsDeleted{Status: messages.Failed, EventId: deleteInfo.EventId},
			fmt.Errorf("can`t delete user: %v; rollback err: %v", err, rbErr),
		), op)
		//TODO: не уверен, что ошибки так заворачиваются
	}

	if _, err := g.requestsRepository.Create(ctx, repositories_transfer.CreateRequestInfo{
		Key: deleteInfo.EventId,
	}, tx); err != nil {
		g.log.Error("can`t create request from db :", err)
		rbErr := tx.Rollback()
		return domain.AddOpInErr(g.createPostDeleteFeedbackEvent(
			ctx,
			messages.PostsDeleted{Status: messages.Failed, EventId: deleteInfo.EventId},
			fmt.Errorf("can`t delete user: %v; rollback err: %v", err, rbErr),
		), op)
	}

	if err := g.createPostDeleteFeedbackEvent(
		ctx,
		messages.PostsDeleted{Status: messages.Success, EventId: deleteInfo.EventId},
		nil,
	); err != nil {
		log.Error("can`t create post delete request from db :", err)
		return domain.AddOpInErr(err, op)
	}
	return nil
}

func (g *PostService) createPostDeleteFeedbackEvent(ctx context.Context, message messages.PostsDeleted, outErr error) error { //TODO: ПОЧЕКАТЬ МОЖНО ЛИ ДАВАТЬ ЗНАЧЕНИЯ ПО УМОЛЧАНИЮ P.S: В РЕПОЗИТОРИЯХ ТОЖЕ САМОЕ
	op := "PostsService.createPostDeleteFeedbackEvent"

	payload, err := json.Marshal(message)
	if err != nil {
		return domain.AddOpInErr(err, op)
	}

	if err := g.eventsRepository.Create(ctx, repositories_transfer.CreateEventInfo{
		EventId:   uuid.New(),
		EventType: amqpclient.PostsDeletedEventKey,
		Payload:   payload,
	}, nil); err != nil {
		return domain.AddOpInErr(err, op)
	}

	return nil
}
