package grpcserver

import (
	"context"
	"fmt"
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	commentservice "github.com/KBcHMFollower/blog_posts_service/internal/services"
	"github.com/google/uuid"
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
	postId, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, fmt.Errorf("error parsing post id: %w", err)
	}

	comments, err := s.commService.GetPostComments(ctx, &services_transfer.GetPostCommentsInfo{
		PostId: postId,
		Size:   req.Size,
		Page:   req.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting post comments: %w", err)
	}

	return &commentsv1.GetPostCommentsResponse{
		Comments:   services_transfer.CommentsArrayProto(comments.Comments),
		TotalCount: int32(comments.TotalCount),
	}, nil
}

func (s *CommentsServer) GetComment(ctx context.Context, req *commentsv1.GetCommentRequest) (*commentsv1.GetCommentResponse, error) {
	commentId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("error parsing post id: %w", err)
	}

	comm, err := s.commService.GetComment(ctx, &services_transfer.GetCommentInfo{
		CommId: commentId,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting comment: %w", err)
	}

	return &commentsv1.GetCommentResponse{
		Comments: comm.Comment.ToProto(), //TODO: COMMENT
	}, nil
}

func (s *CommentsServer) DeleteComment(ctx context.Context, req *commentsv1.DeleteCommentRequest) (*commentsv1.DeleteCommentResponse, error) {
	commentId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("error parsing post id: %w", err)
	}

	if err := s.commService.DeleteComment(ctx, &services_transfer.DeleteCommentInfo{
		CommId: commentId,
	}); err != nil {
		return &commentsv1.DeleteCommentResponse{
			IsDeleted: false,
		}, fmt.Errorf("error deleting comment: %w", err)
	}

	return &commentsv1.DeleteCommentResponse{
		IsDeleted: true,
	}, nil
}

func (s *CommentsServer) UpdateComment(ctx context.Context, req *commentsv1.UpdateCommentRequest) (*commentsv1.UpdateCommentResponse, error) {
	commentId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("error parsing post id: %w", err)
	}

	var fields = make([]services_transfer.CommUpdateFieldInfo, 0, len(req.UpdateData))
	for _, item := range req.UpdateData {
		fields = append(fields, services_transfer.CommUpdateFieldInfo{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	comm, err := s.commService.UpdateComment(ctx, &services_transfer.UpdateCommentInfo{
		CommId:       commentId,
		UpdateFields: fields,
	})
	if err != nil {
		return nil, fmt.Errorf("error updating comment: %w", err)
	}

	return &commentsv1.UpdateCommentResponse{
		Comment: comm.Comment.ToProto(),
		Id:      commentId.String(),
	}, nil
}

func (s *CommentsServer) CreateComment(ctx context.Context, req *commentsv1.CreateCommentRequest) (*commentsv1.CreateCommentResponse, error) {
	postId, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, fmt.Errorf("error parsing post id: %w", err)
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("error parsing post id: %w", err)
	}

	comm, err := s.commService.CreateComment(ctx, &services_transfer.CreateCommentInfo{
		UserId:  userId,
		PostId:  postId,
		Content: req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	return &commentsv1.CreateCommentResponse{
		Id:      comm.CommId.String(),
		Comment: comm.Comment.ToProto(),
	}, nil
}
