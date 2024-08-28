package services

import (
	"context"
	"fmt"
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	services_dep "github.com/KBcHMFollower/blog_posts_service/internal/services/interfaces/dep"
	"log/slog"
)

type CommentsStore interface {
	services_dep.CommentsCreator
	services_dep.CommentsGetter
	services_dep.CommentsDeleter
	services_dep.CommentsUpdater
}

type CommentsService struct {
	commRep CommentsStore
	log     *slog.Logger
}

func NewCommentService(commReP CommentsStore, log *slog.Logger) *CommentsService {
	return &CommentsService{
		commRep: commReP,
		log:     log,
	}
}

func (s *CommentsService) GetPostComments(ctx context.Context, getInfo *services_transfer.GetPostCommentsInfo) (*services_transfer.GetPostCommentsResult, error) {
	op := "CommentsService.GetPostComments"
	log := s.log.With(
		slog.String("op", op),
	)

	comments, err := s.commRep.GetPostComments(ctx, getInfo.PostId, uint64(getInfo.Size), uint64(getInfo.Page))
	if err != nil {
		log.Error(fmt.Sprintf("failed to get comments for post %d: %v", getInfo.PostId, err))
		return nil, domain.AddOpInErr(err, op)
	}
	total, err := s.commRep.GetPostCommentsCount(ctx, getInfo.PostId)
	if err != nil {
		log.Error(fmt.Sprintf("failed to get comments count for post %d: %v", getInfo.PostId, err))
		return nil, domain.AddOpInErr(err, op)
	}

	resComms := make([]*commentsv1.Comment, 0)

	for _, item := range comments {
		resPost := item.ConvertToProto()
		resComms = append(resComms, resPost)
	}
	return &services_transfer.GetPostCommentsResult{
		TotalCount: int32(total),
		Comments:   services_transfer.ConvertCommentsArrayFromModels(comments),
	}, nil
}

func (s *CommentsService) GetComment(ctx context.Context, getInfo *services_transfer.GetCommentInfo) (*services_transfer.GetCommentResult, error) {
	op := "CommentsService.GetComment"
	log := s.log.With(
		slog.String("op", op),
	)

	comment, err := s.commRep.GetComment(ctx, getInfo.CommId)
	if err != nil {
		log.Error(fmt.Sprintf("failed to get comment for comment %d: %v", getInfo.CommId, err))
		return nil, domain.AddOpInErr(err, op)
	}

	return &services_transfer.GetCommentResult{
		Comment: services_transfer.ConvertCommentFromModel(comment),
	}, nil
}

func (s *CommentsService) DeleteComment(ctx context.Context, deleteInfo *services_transfer.DeleteCommentInfo) error {
	op := "PostService.DeletePost"

	log := s.log.With(
		slog.String("op", op),
	)

	if err := s.commRep.DeleteComment(ctx, deleteInfo.CommId); err != nil {
		log.Error(fmt.Sprintf("failed to delete comment %d: %v", deleteInfo.CommId, err))
		return domain.AddOpInErr(err, op)
	}

	return nil
}

func (s *CommentsService) UpdateComment(ctx context.Context, updateInfo *services_transfer.UpdateCommentInfo) (*services_transfer.UpdateCommentResult, error) {
	op := "CommentsService.UpdateComment"
	log := s.log.With(
		slog.String("op", op),
	)

	updateItems := make([]*repositories_transfer.CommentUpdateFieldInfo, 0)

	for _, item := range updateInfo.UpdateFields {
		updateItems = append(updateItems, &repositories_transfer.CommentUpdateFieldInfo{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	if err := s.commRep.UpdateComment(ctx, repositories_transfer.UpdateCommentInfo{
		Id:         updateInfo.CommId,
		UpdateData: updateItems,
	}); err != nil {
		log.Error(fmt.Sprintf("failed to update comment %d: %v", updateInfo.CommId, err))
		return nil, domain.AddOpInErr(err, op)
	}
	comm, err := s.commRep.GetComment(ctx, updateInfo.CommId)
	if err != nil {
		log.Error(fmt.Sprintf("failed to get comment %d: %v", updateInfo.CommId, err))
		return nil, domain.AddOpInErr(err, op)
	}

	return &services_transfer.UpdateCommentResult{
		CommId:  comm.Id,
		Comment: services_transfer.ConvertCommentFromModel(comm),
	}, nil
}

func (s *CommentsService) CreateComment(ctx context.Context, createInfo *services_transfer.CreateCommentInfo) (*services_transfer.CreateCommentResult, error) {
	op := "CommentsService.CreateComment"
	log := s.log.With(
		slog.String("op", op),
	)

	commId, err := s.commRep.CreateComment(ctx, repositories_transfer.CreateCommentInfo{
		PostId:  createInfo.PostId,
		UserId:  createInfo.UserId,
		Content: createInfo.Content, //TODO: ВРОДЕ ЕЩЕ КАРТИНКИ МОЖНО
	})
	if err != nil {
		log.Error(fmt.Sprintf("failed to create comment %d: %v", createInfo.PostId, err))
		return nil, domain.AddOpInErr(err, op)
	}
	comm, err := s.commRep.GetComment(ctx, commId)

	return &services_transfer.CreateCommentResult{
		CommId:  commId,
		Comment: services_transfer.ConvertCommentFromModel(comm),
	}, nil
}
