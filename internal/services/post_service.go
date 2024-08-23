package services

import (
	"context"
	"fmt"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	services_dep "github.com/KBcHMFollower/blog_posts_service/internal/services/interfaces/dep"
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

type PostService struct {
	postRepository     PostsStore
	requestsRepository RequestStore
	log                *slog.Logger
}

func NewPostsService(postRepository PostsStore, log *slog.Logger) *PostService {
	return &PostService{
		postRepository: postRepository,
		log:            log,
	}
}

func (g *PostService) GetUserPosts(ctx context.Context, getInfo *services_transfer.GetUserPostsInfo) (*services_transfer.GetUserPostsResult, error) {
	op := "PostService.GetUserPosts"

	log := g.log.With(
		slog.String("op", op),
	)

	posts, totalCount, err := g.postRepository.GetPostsByUserId(ctx, getInfo.UserId, uint64(getInfo.Size), uint64(getInfo.Page))
	if err != nil {
		log.Error("repository error", err)
		return nil, err
	}

	return &services_transfer.GetUserPostsResult{
		Posts:      services_transfer.ConvertPostsArrayFromModel(posts),
		TotalCount: int32(totalCount),
	}, nil
}

func (g *PostService) GetPost(ctx context.Context, getInfo *services_transfer.GetPostInfo) (*services_transfer.GetPostResult, error) {
	op := "PostService.GetPost"

	log := g.log.With(
		slog.String("op", op),
	)

	post, err := g.postRepository.GetPost(ctx, getInfo.PostId)
	if err != nil {
		log.Error("can`t get user from db :", err)
		return nil, err
	}

	return &services_transfer.GetPostResult{
		Post: services_transfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) DeletePost(ctx context.Context, deleteInfo *services_transfer.DeletePostInfo) error {
	op := "PostService.DeletePost"

	log := g.log.With(
		slog.String("op", op),
	)

	if _, err := g.postRepository.DeletePost(ctx, deleteInfo.PostId); err != nil {
		log.Error("repository error", err)
		return err
	}

	return nil
}

func (g *PostService) CreatePost(ctx context.Context, createInfo *services_transfer.CreatePostInfo) (*services_transfer.CreatePostResult, error) {
	op := "PostService.CreatePost"

	log := g.log.With(
		slog.String("op", op),
	)

	postId, post, err := g.postRepository.CreatePost(ctx, repositories_transfer.CreatePostInfo{
		UserId:        createInfo.UserId,
		Title:         createInfo.Title,
		TextContent:   createInfo.TextContent,
		ImagesContent: createInfo.ImagesContent,
	})
	if err != nil {
		log.Error("can`t create user from db :", err)
		return nil, err
	}

	fmt.Println(postId, post)

	return &services_transfer.CreatePostResult{
		PostId: postId,
		Post:   services_transfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) UpdatePost(ctx context.Context, updateInfo *services_transfer.UpdatePostInfo) (*services_transfer.UpdatePostResult, error) {
	op := "PostService.CreatePost"

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

	post, err := g.postRepository.UpdatePost(ctx, repositories_transfer.UpdateCommentInfo{
		Id:         updateInfo.PostId,
		UpdateData: updateItems,
	})
	if err != nil {
		log.Error("can`t update user from db :", err)
		return nil, err
	}

	return &services_transfer.UpdatePostResult{
		PostId: post.Id,
		Post:   services_transfer.ConvertPostFromModel(post),
	}, nil
}

func (g *PostService) DeleteUserPosts(ctx context.Context, deleteInfo services_transfer.DeleteUserPostInfo) error {
	req, err := g.requestsRepository.Get() //TODO
}
