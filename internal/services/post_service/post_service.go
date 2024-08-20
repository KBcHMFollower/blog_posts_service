package post_service

import (
	"context"
	"fmt"
	postsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/posts"
	"log/slog"

	"github.com/KBcHMFollower/blog_posts_service/internal/repository"
	"github.com/google/uuid"
)

type PostService struct {
	postRepository repository.PostsStore
	log            *slog.Logger
}

func New(postRepository repository.PostsStore, log *slog.Logger) *PostService {
	return &PostService{
		postRepository: postRepository,
		log:            log,
	}
}

func (g *PostService) GetUserPosts(ctx context.Context, req *postsv1.GetUserPostsRequest) (*postsv1.GetUserPostsResponse, error) {
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

	resPosts := make([]*postsv1.Post, 0)

	for _, item := range posts {
		resPost := item.ConvertToProto()
		resPosts = append(resPosts, resPost)
	}
	return &postsv1.GetUserPostsResponse{
		Posts:      resPosts,
		TotalCount: int32(totalCount),
	}, nil
}

func (g *PostService) GetPost(ctx context.Context, req *postsv1.GetPostRequest) (*postsv1.GetPostResponse, error) {
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

	return &postsv1.GetPostResponse{
		Posts: post.ConvertToProto(),
	}, nil
}

func (g *PostService) DeletePost(ctx context.Context, req *postsv1.DeletePostRequest) (*postsv1.DeletePostResponse, error) {
	op := "PostService.DeletePost"

	log := g.log.With(
		slog.String("op", op),
	)

	postUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	_, err = g.postRepository.DeletePost(ctx, postUUID)
	if err != nil {
		log.Error("can`t delete user from db :", err)
		return &postsv1.DeletePostResponse{
			IsDeleted: false,
		}, err
	}

	return &postsv1.DeletePostResponse{
		IsDeleted: true,
	}, nil
}

func (g *PostService) CreatePost(ctx context.Context, req *postsv1.CreatePostRequest) (*postsv1.CreatePostResponse, error) {
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

	fmt.Println(postId, post)

	return &postsv1.CreatePostResponse{
		Id:   postId.String(),
		Post: post.ConvertToProto(),
	}, nil
}

func (g *PostService) UpdatePost(ctx context.Context, req *postsv1.UpdatePostRequest) (*postsv1.UpdatePostResponse, error) {
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

	post, err := g.postRepository.UpdatePost(ctx, repository.UpdateData{
		Id:         postUUID,
		UpdateData: updateItems,
	})
	if err != nil {
		log.Error("can`t update user from db :", err)
		return nil, err
	}

	return &postsv1.UpdatePostResponse{
		Id:   post.Id.String(),
		Post: post.ConvertToProto(),
	}, nil
}

func (g *PostService) DeleteUserPosts(ctx context.Context, userId uuid.UUID) error {
	err := g.postRepository.DeleteUserPosts(ctx, userId)
	return err
}
