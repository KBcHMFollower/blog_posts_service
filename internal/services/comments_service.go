package services

import (
	"context"
	"fmt"
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	ctxerrors "github.com/KBcHMFollower/blog_posts_service/internal/domain/errors"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	services_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/services"
	"github.com/KBcHMFollower/blog_posts_service/internal/logger"
	services_dep "github.com/KBcHMFollower/blog_posts_service/internal/services/interfaces/dep"
)

type CommentsStore interface {
	services_dep.CommentsCreator
	services_dep.CommentsGetter
	services_dep.CommentsDeleter
	services_dep.CommentsUpdater
}

type CommentsService struct {
	commRep CommentsStore
	log     logger.Logger
}

func NewCommentService(commReP CommentsStore, log logger.Logger) *CommentsService {
	return &CommentsService{
		commRep: commReP,
		log:     log,
	}
}

func (s *CommentsService) GetPostComments(ctx context.Context, getInfo *services_transfer.GetPostCommentsInfo) (*services_transfer.GetPostCommentsResult, error) {
	comments, err := s.commRep.Comments(ctx, repositories_transfer.GetCommentsInfo{
		Size: getInfo.Size,
		Page: getInfo.Page,
		Condition: map[repositories_transfer.CommentConditionTarget]interface{}{
			repositories_transfer.CommentPostIdConditionTarget: getInfo.PostId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get comments from repository", err))
	}

	total, err := s.commRep.Count(ctx, repositories_transfer.GetCommentsCountInfo{
		Condition: map[repositories_transfer.CommentConditionTarget]interface{}{
			repositories_transfer.CommentPostIdConditionTarget: getInfo.PostId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get comments count from repository", err))
	}

	resComms := make([]*commentsv1.Comment, 0)

	for _, item := range comments {
		resPost := item.ConvertToProto()
		resComms = append(resComms, resPost)
	}
	return &services_transfer.GetPostCommentsResult{
		TotalCount: total,
		Comments:   services_transfer.ConvertCommentsArrayFromModels(comments),
	}, nil
}

func (s *CommentsService) GetComment(ctx context.Context, getInfo *services_transfer.GetCommentInfo) (*services_transfer.GetCommentResult, error) {
	comment, err := s.commRep.Comment(ctx, repositories_transfer.GetCommentInfo{
		Condition: map[repositories_transfer.CommentConditionTarget]interface{}{
			repositories_transfer.CommentIdConditionTarget: getInfo.CommId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get comment from repository", err))
	}

	return &services_transfer.GetCommentResult{
		Comment: services_transfer.ConvertCommentFromModel(comment),
	}, nil
}

func (s *CommentsService) DeleteComment(ctx context.Context, deleteInfo *services_transfer.DeleteCommentInfo) error {
	logger.UpdateLoggerCtx(ctx, "comment-id", deleteInfo.CommId)
	s.log.InfoContext(ctx, "try to delete comment")

	if err := s.commRep.Delete(ctx, repositories_transfer.DeleteCommentInfo{
		Condition: map[repositories_transfer.CommentConditionTarget]interface{}{
			repositories_transfer.CommentIdConditionTarget: deleteInfo.CommId,
		},
	}); err != nil {
		return ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant delete comments from repository", err))
	}

	s.log.InfoContext(ctx, "comment is deleted")

	return nil
}

func (s *CommentsService) UpdateComment(ctx context.Context, updateInfo *services_transfer.UpdateCommentInfo) (*services_transfer.UpdateCommentResult, error) {
	logger.UpdateLoggerCtx(ctx, "comment-id", updateInfo.CommId)
	logger.UpdateLoggerCtx(ctx, "update-data", updateInfo.UpdateFields)
	s.log.InfoContext(ctx, "try to update comment")

	if err := s.commRep.Update(ctx, repositories_transfer.UpdateCommentInfo{
		Condition: map[repositories_transfer.CommentConditionTarget]interface{}{
			repositories_transfer.CommentIdConditionTarget: updateInfo.CommId,
		},
		UpdateData: updateInfo.UpdateFields,
	}); err != nil {
		s.log.Error(fmt.Sprintf("failed to update comment %d: %v", updateInfo.CommId, err))
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant update comments in repository", err))
	}

	comm, err := s.commRep.Comment(ctx, repositories_transfer.GetCommentInfo{
		Condition: map[repositories_transfer.CommentConditionTarget]interface{}{
			repositories_transfer.CommentIdConditionTarget: updateInfo.CommId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get comment from repository", err))
	}

	s.log.InfoContext(ctx, "comment is updated")

	return &services_transfer.UpdateCommentResult{
		CommId:  comm.Id,
		Comment: services_transfer.ConvertCommentFromModel(comm),
	}, nil
}

func (s *CommentsService) CreateComment(ctx context.Context, createInfo *services_transfer.CreateCommentInfo) (*services_transfer.CreateCommentResult, error) {
	logger.UpdateLoggerCtx(ctx, "create-info", createInfo)
	s.log.InfoContext(ctx, "try to create comment")

	commId, err := s.commRep.Create(ctx, repositories_transfer.CreateCommentInfo{
		PostId:  createInfo.PostId,
		UserId:  createInfo.UserId,
		Content: createInfo.Content,
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant create comments in repository", err))
	}

	comm, err := s.commRep.Comment(ctx, repositories_transfer.GetCommentInfo{
		Condition: map[repositories_transfer.CommentConditionTarget]interface{}{
			repositories_transfer.CommentIdConditionTarget: commId,
		},
	})
	if err != nil {
		return nil, ctxerrors.WrapCtx(ctx, ctxerrors.Wrap("cant get comment from repository", err))
	}

	s.log.InfoContext(ctx, "comment is created")

	return &services_transfer.CreateCommentResult{
		CommId:  commId,
		Comment: services_transfer.ConvertCommentFromModel(comm),
	}, nil
}
