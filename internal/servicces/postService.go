package postService

import (
	"context"
	"log/slog"

	ssov1 "github.com/KBcHMFollower/test_plate_user_service/internal/api/protos/gen"
	"github.com/KBcHMFollower/test_plate_user_service/internal/repository"
	"github.com/google/uuid"
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

func (g *PostService) GetUserPosts(ctx context.Context, req *ssov1.GetUserPostsRequest) (*ssov1.GetUserPostsResponse, error) {
	op := "PostService.GetUserPosts"

	log := g.log.With(
		slog.String("op", op),
	)

	userUUID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Error("can`t parse user_id from uuid", err)
		return nil, err
	}

	posts, totalCount, err := g.postRepository.GetPostsByUserId(ctx, userUUID, uint64(req.GetSize()), uint64(req.GetPage()))
	if err != nil {
		log.Error("repository error", err)
		return nil, err
	}

	resPosts := make([]*ssov1.Post, 0)

	for _, item := range posts {
		resPost := item.ConvertToProto()
		resPosts = append(resPosts, resPost)
	}
	return &ssov1.GetUserPostsResponse{
		Posts:      resPosts,
		TotalCount: int32(totalCount),
	}, nil
}

func (g *PostService) GetPost(ctx context.Context, req *ssov1.GetPostRequest) (*ssov1.GetPostResponse, error) {
	op := "PostService.GetPost"

	log := g.log.With(
		slog.String("op", op),
	)

	postUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	post, err := g.postRepository.GetPost(ctx, postUUID)
	if err != nil {
		log.Error("can`t get user from db :", err)
		return nil, err
	}

	return &ssov1.GetPostResponse{
		Posts: post.ConvertToProto(),
	}, nil
}

func (g *PostService) DeletePost(ctx context.Context, req *ssov1.DeletePostRequest) (*ssov1.DeletePostResponse, error) {
	op := "PostService.DeletePost"

	log := g.log.With(
		slog.String("op", op),
	)

	postUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	err = g.postRepository.DeletePost(ctx, postUUID)
	if err != nil {
		log.Error("can`t delete user from db :", err)
		return &ssov1.DeletePostResponse{
			IsDeleted: false,
		}, err
	}

	return &ssov1.DeletePostResponse{
		IsDeleted: true,
	}, nil
}

func (g *PostService) CreatePost(ctx context.Context, req *ssov1.CreatePostRequest) (*ssov1.CreatePostResponse, error) {
	op := "PostService.CreatePost"

	log := g.log.With(
		slog.String("op", op),
	)

	userUUID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	postId, post, err := g.postRepository.CreatePost(ctx, repository.CreatePostData{
		User_id:       userUUID,
		Title:         req.GetTitle(),
		TextContent:   req.GetTextContent(),
		ImagesContent: req.ImagesContent,
	})
	if err != nil {
		log.Error("can`t create user from db :", err)
		return nil, err
	}

	return &ssov1.CreatePostResponse{
		Id:   postId.String(),
		Post: post.ConvertToProto(),
	}, nil
}

func (g *PostService) UpdatePost(ctx context.Context, req *ssov1.UpdatePostRequest) (*ssov1.UpdatePostResponse, error) {
	op := "PostService.CreatePost"

	log := g.log.With(
		slog.String("op", op),
	)

	postUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	updateItems := make([]*repository.UpdateItem, 0)

	for _, item := range req.UpdateData {
		updateItems = append(updateItems, &repository.UpdateItem{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	post, err := g.postRepository.UpdatePost(ctx, repository.UpdatePostData{
		Id:         postUUID,
		UpdateData: updateItems,
	})
	if err != nil {
		log.Error("can`t update user from db :", err)
		return nil, err
	}

	return &ssov1.UpdatePostResponse{
		Id:   post.Id.String(),
		Post: post.ConvertToProto(),
	}, nil
}
