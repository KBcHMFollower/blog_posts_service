package grpcserver

import (
	"context"
	postService "github.com/KBcHMFollower/test_plate_blog_service/internal/services/post_service"

	postsv1 "github.com/KBcHMFollower/test_plate_blog_service/api/protos/gen/posts"
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
	return g.postService.GetUserPosts(ctx, req)
}

func (g *PostsServer) GetPost(ctx context.Context, req *postsv1.GetPostRequest) (*postsv1.GetPostResponse, error) {
	return g.postService.GetPost(ctx, req)
}

func (g *PostsServer) DeletePost(ctx context.Context, req *postsv1.DeletePostRequest) (*postsv1.DeletePostResponse, error) {
	return g.postService.DeletePost(ctx, req)
}

func (g *PostsServer) UpdatePost(ctx context.Context, req *postsv1.UpdatePostRequest) (*postsv1.UpdatePostResponse, error) {
	return g.postService.UpdatePost(ctx, req)
}

func (g *PostsServer) CreatePost(ctx context.Context, req *postsv1.CreatePostRequest) (*postsv1.CreatePostResponse, error) {
	return g.postService.CreatePost(ctx, req)
}
