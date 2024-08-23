package grpcserver

import (
	"context"
	"fmt"
	postsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/posts"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	postService "github.com/KBcHMFollower/blog_posts_service/internal/services"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type PostsServer struct {
	postsv1.UnimplementedPostServiceServer
	postService *postService.PostService
}

func RegisterPostServer(server *grpc.Server, postService *postService.PostService) {
	postsv1.RegisterPostServiceServer(server, &PostsServer{postService: postService})
}

func (g *PostsServer) GetUserPosts(ctx context.Context, req *postsv1.GetUserPostsRequest) (*postsv1.GetUserPostsResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %s", req.UserId)
	}

	posts, err := g.postService.GetUserPosts(ctx, &services_transfer.GetUserPostsInfo{
		UserId: userId,
		Size:   req.Size,
		Page:   req.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting posts for user id %s", req.UserId)
	}

	return &postsv1.GetUserPostsResponse{
		Posts:      services_transfer.ConvertPostArrayToProto(posts.Posts),
		TotalCount: int32(posts.TotalCount),
	}, nil
}

func (g *PostsServer) GetPost(ctx context.Context, req *postsv1.GetPostRequest) (*postsv1.GetPostResponse, error) {
	postId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %s", req.Id)
	}

	Post, err := g.postService.GetPost(ctx, &services_transfer.GetPostInfo{
		PostId: postId,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting post %s", req.Id)
	}

	return &postsv1.GetPostResponse{
		Posts: Post.Post.ToProto(),
	}, nil //TODO: POST
}

func (g *PostsServer) DeletePost(ctx context.Context, req *postsv1.DeletePostRequest) (*postsv1.DeletePostResponse, error) {
	postId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %s", req.Id)
	}

	if err := g.postService.DeletePost(ctx, &services_transfer.DeletePostInfo{
		PostId: postId,
	}); err != nil {
		return &postsv1.DeletePostResponse{
			IsDeleted: false,
		}, fmt.Errorf("error deleting post %s", req.Id)
	}

	return &postsv1.DeletePostResponse{
		IsDeleted: true,
	}, nil
}

func (g *PostsServer) UpdatePost(ctx context.Context, req *postsv1.UpdatePostRequest) (*postsv1.UpdatePostResponse, error) {
	postId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %s", req.Id)
	}

	var fields []services_transfer.UpdateUserFieldInfo = make([]services_transfer.UpdateUserFieldInfo, 0, len(req.UpdateData))
	for _, item := range req.UpdateData {
		fields = append(fields, services_transfer.UpdateUserFieldInfo{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	post, err := g.postService.UpdatePost(ctx, &services_transfer.UpdatePostInfo{
		PostId: postId,
		Fields: fields,
	})
	if err != nil {
		return nil, fmt.Errorf("error updating post %s", req.Id)
	}

	return &postsv1.UpdatePostResponse{
		Post: post.Post.ToProto(),
		Id:   post.PostId.String(),
	}, nil
}

func (g *PostsServer) CreatePost(ctx context.Context, req *postsv1.CreatePostRequest) (*postsv1.CreatePostResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %s", req.UserId)
	}

	post, err := g.postService.CreatePost(ctx, &services_transfer.CreatePostInfo{
		UserId:        userId,
		Title:         req.Title,
		TextContent:   req.TextContent,
		ImagesContent: req.ImagesContent,
	})
	if err != nil {
		return nil, fmt.Errorf("error in creating post: %w", err)
	}

	return &postsv1.CreatePostResponse{
		Id:   post.PostId.String(),
		Post: post.Post.ToProto(),
	}, nil
}
