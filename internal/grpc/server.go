package grpcserver

import (
	"context"
	"fmt"

	ssov1 "github.com/KBcHMFollower/test_plate_user_service/internal/api/protos/gen"
	postService "github.com/KBcHMFollower/test_plate_user_service/internal/servicces"
	"google.golang.org/grpc"
)

type GPRCApp struct {
	ssov1.UnimplementedBlogsServer
	postService postService.PostService
}

func Register(server *grpc.Server, postService postService.PostService) {
	ssov1.RegisterBlogsServer(server, &GPRCApp{postService: postService})
}

func (g *GPRCApp) GetUserPosts(ctx context.Context, req *ssov1.GetUserPostsRequest) (*ssov1.GetUserPostsResponse, error) {
	return nil, fmt.Errorf("method GetUserPosts not implemented")
}

func (g *GPRCApp) GetPost(ctx context.Context, req *ssov1.GetPostRequest) (*ssov1.GetPostResponse, error) {
	return nil, fmt.Errorf("method GetUserPosts not implemented")
}

func (g *GPRCApp) DeletePost(ctx context.Context, req *ssov1.DeletePostRequest) (*ssov1.DeletePostResponse, error) {
	return nil, fmt.Errorf("method GetUserPosts not implemented")
}

func (g *GPRCApp) UpdatePost(ctx context.Context, req *ssov1.UpdatePostRequest) (*ssov1.UpdatePostResponse, error) {
	return nil, fmt.Errorf("method GetUserPosts not implemented")
}

func (g *GPRCApp) CreatePost(ctx context.Context, req *ssov1.CreatePostRequest) (*ssov1.CreatePostResponse, error) {
	return nil, fmt.Errorf("method GetUserPosts not implemented")
}
