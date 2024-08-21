package grpcserver

import (
	"context"
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	commentservice "github.com/KBcHMFollower/blog_posts_service/internal/services"
	"google.golang.org/grpc"
)

type CommentsServer struct {
	commentsv1.UnimplementedCommentServiceServer
	commService *commentservice.CommentsService
}

func RegisterCommentsServer(server *grpc.Server, commService *commentservice.CommentsService) {
	commentsv1.RegisterCommentServiceServer(server, &CommentsServer{commService: commService})
}

func (s *CommentsServer) GetPostComments(ctx context.Context, req *commentsv1.GetPostCommentsRequest) (*commentsv1.GetPostCommentsResponse, error) {
	return s.commService.GetPostComments(ctx, req)
}

func (s *CommentsServer) GetComment(ctx context.Context, req *commentsv1.GetCommentRequest) (*commentsv1.GetCommentResponse, error) {
	return s.commService.GetComment(ctx, req)
}

func (s *CommentsServer) DeleteComment(ctx context.Context, req *commentsv1.DeleteCommentRequest) (*commentsv1.DeleteCommentResponse, error) {
	return s.commService.DeleteComment(ctx, req)
}

func (s *CommentsServer) UpdateComment(ctx context.Context, req *commentsv1.UpdateCommentRequest) (*commentsv1.UpdateCommentResponse, error) {
	return s.commService.UpdateComment(ctx, req)
}

func (s *CommentsServer) CreateComment(ctx context.Context, req *commentsv1.CreateCommentRequest) (*commentsv1.CreateCommentResponse, error) {
	return s.commService.CreateComment(ctx, req)
}
