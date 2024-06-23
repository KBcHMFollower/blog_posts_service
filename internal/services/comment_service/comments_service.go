package commentservice

import (
	"context"
	commentsv1 "github.com/KBcHMFollower/test_plate_blog_service/api/protos/gen/comments"
	"github.com/KBcHMFollower/test_plate_blog_service/internal/repository"
	"github.com/google/uuid"
	"log/slog"
)

type CommentsService struct {
	commRep repository.CommentStore
	log     *slog.Logger
}

func New(commReP repository.CommentStore, log *slog.Logger) *CommentsService {
	return &CommentsService{
		commRep: commReP,
		log:     log,
	}
}

func (s *CommentsService) GetPostComments(ctx context.Context, req *commentsv1.GetPostCommentsRequest) (*commentsv1.GetPostCommentsResponse, error) {
	op := "CommentsService.GetPostComments"

	log := s.log.With(
		slog.String("op", op),
	)

	postUUID, err := uuid.Parse(req.GetPostId())
	if err != nil {
		log.Error("can`t parse user_id from uuid", err)
		return nil, err
	}

	comments, totalCount, err := s.commRep.GetPostComments(ctx, postUUID, uint64(req.GetSize()), uint64(req.GetPage()))
	if err != nil {
		log.Error("repository error", err)
		return nil, err
	}

	resComms := make([]*commentsv1.Comment, 0)

	for _, item := range comments {
		resPost := item.ConvertToProto()
		resComms = append(resComms, resPost)
	}
	return &commentsv1.GetPostCommentsResponse{
		Comments:   resComms,
		TotalCount: int32(totalCount),
	}, nil
}

func (s *CommentsService) GetComment(ctx context.Context, req *commentsv1.GetCommentRequest) (*commentsv1.GetCommentResponse, error) {
	op := "CommentsService.GetComment"

	log := s.log.With(
		slog.String("op", op),
	)

	commentUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	comment, err := s.commRep.GetComment(ctx, commentUUID)
	if err != nil {
		log.Error("can`t get user from db :", err)
		return nil, err
	}

	return &commentsv1.GetCommentResponse{
		Comments: comment.ConvertToProto(),
	}, nil
}

func (s *CommentsService) DeleteComment(ctx context.Context, req *commentsv1.DeleteCommentRequest) (*commentsv1.DeleteCommentResponse, error) {
	op := "PostService.DeletePost"

	log := s.log.With(
		slog.String("op", op),
	)

	commentUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	_, err = s.commRep.DeleteComment(ctx, commentUUID)
	if err != nil {
		log.Error("can`t delete user from db :", err)
		return &commentsv1.DeleteCommentResponse{
			IsDeleted: false,
		}, err
	}

	return &commentsv1.DeleteCommentResponse{
		IsDeleted: true,
	}, nil
}

func (s *CommentsService) UpdateComment(ctx context.Context, req *commentsv1.UpdateCommentRequest) (*commentsv1.UpdateCommentResponse, error) {
	op := "CommentsService.UpdateComment"

	log := s.log.With(
		slog.String("op", op),
	)

	commUUID, err := uuid.Parse(req.GetId())
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

	comm, err := s.commRep.UpdateComment(ctx, repository.UpdateData{
		Id:         commUUID,
		UpdateData: updateItems,
	})
	if err != nil {
		log.Error("can`t update user from db :", err)
		return nil, err
	}

	return &commentsv1.UpdateCommentResponse{
		Id:      comm.Id.String(),
		Comment: comm.ConvertToProto(),
	}, nil
}

func (s *CommentsService) CreateComment(ctx context.Context, req *commentsv1.CreateCommentRequest) (*commentsv1.CreateCommentResponse, error) {
	op := "CommentsService.CreateComment"

	log := s.log.With(
		slog.String("op", op),
	)

	userUUID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Error("can`t parse user_id from uuid :", err)
		return nil, err
	}

	postUUID, err := uuid.Parse(req.GetPostId())
	if err != nil {
		log.Error("can`t parse post_id from uuid :", err)
		return nil, err
	}

	commId, comm, err := s.commRep.CreateComment(ctx, repository.CreateCommentData{
		PostId:  postUUID,
		UserId:  userUUID,
		Content: req.GetContent(),
	})
	if err != nil {
		log.Error("can`t create user from db :", err)
		return nil, err
	}

	return &commentsv1.CreateCommentResponse{
		Id:      commId.String(),
		Comment: comm.ConvertToProto(),
	}, nil
}
