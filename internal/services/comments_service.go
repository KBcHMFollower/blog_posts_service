package services

import (
	"context"
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
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

	comments, totalCount, err := s.commRep.GetPostComments(ctx, getInfo.PostId, uint64(getInfo.Size), uint64(getInfo.Page))
	if err != nil {
		log.Error("repository error", err)
		return nil, err
	}

	resComms := make([]*commentsv1.Comment, 0)

	for _, item := range comments {
		resPost := item.ConvertToProto()
		resComms = append(resComms, resPost)
	}
	return &services_transfer.GetPostCommentsResult{
		TotalCount: int32(totalCount),
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
		log.Error("can`t get user from db :", err)
		return nil, err
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

	if _, err := s.commRep.DeleteComment(ctx, deleteInfo.CommId); err != nil {
		log.Error("can`t delete user from db :", err)
		return err
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

	comm, err := s.commRep.UpdateComment(ctx, repositories_transfer.UpdateCommentInfo{
		Id:         updateInfo.CommId,
		UpdateData: updateItems,
	})
	if err != nil {
		log.Error("can`t update user from db :", err)
		return nil, err
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

	commId, comm, err := s.commRep.CreateComment(ctx, repositories_transfer.CreateCommentInfo{
		PostId:  createInfo.PostId,
		UserId:  createInfo.UserId,
		Content: createInfo.Content, //TODO: ВРОДЕ ЕЩЕ КАРТИНКИ МОЖНО
	})
	if err != nil {
		log.Error("can`t create user from db :", err)
		return nil, err
	}

	return &services_transfer.CreateCommentResult{
		CommId:  commId,
		Comment: services_transfer.ConvertCommentFromModel(comm),
	}, nil
}
