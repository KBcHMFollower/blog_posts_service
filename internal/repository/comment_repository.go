package repository

import (
	"context"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	rep_utils "github.com/KBcHMFollower/blog_posts_service/internal/repository/lib"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	commentsIdCol      = "id"
	commentsUserIdCol  = "user_id"
	commentsPostIdCol  = "post_id"
	commentsContentCol = "content"
	commentsAllCol     = "*"
	commentsSqlCount   = "COUNT(*)"
)

type CommentRepository struct {
	db       database.DBWrapper
	qBuilder squirrel.StatementBuilderType
}

func NewCommentRepository(db database.DBWrapper) *CommentRepository {
	return &CommentRepository{
		db:       db,
		qBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, createData repositories_transfer.CreateCommentInfo) (uuid.UUID, error) {
	op := "CommentRepository.createComment"
	comment := models.CreateComment(createData.PostId, createData.UserId, createData.Content)

	query := r.qBuilder.
		Insert(database.CommentsTable).
		SetMap(map[string]interface{}{
			commentsIdCol:      comment.Id,
			commentsUserIdCol:  comment.UserId,
			commentsPostIdCol:  comment.PostId,
			commentsContentCol: comment.Content,
		}).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, rep_utils.GenerateSqlErr(err, op)
	}

	var insertId uuid.UUID
	if err := r.db.GetContext(ctx, &insertId, sql, args...); err != nil {
		return uuid.Nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return insertId, nil
}

func (r *CommentRepository) GetComment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error) {
	op := "CommentRepository.getComment"
	query := r.qBuilder.
		Select(commentsAllCol).
		From(database.CommentsTable).
		Where(squirrel.Eq{commentsIdCol: commentId})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.GenerateSqlErr(err, op)
	}

	var comment models.Comment

	if err := r.db.GetContext(ctx, &comment, sqlStr, args...); err != nil {
		return nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return &comment, nil
}

func (r *CommentRepository) GetPostCommentsCount(ctx context.Context, postId uuid.UUID) (uint, error) {
	op := "CommentRepository.getPostCommentsCount"

	query := r.qBuilder.
		Select(commentsSqlCount).
		From(database.CommentsTable).
		Where(squirrel.Eq{commentsPostIdCol: postId})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, rep_utils.GenerateSqlErr(err, op)
	}

	var count uint
	if err := r.db.GetContext(ctx, &count, sqlStr, args...); err != nil {
		return 0, rep_utils.ExecuteSqlErr(err, op)
	}

	return count, nil
}

func (r *CommentRepository) GetPostComments(ctx context.Context, postId uuid.UUID, size uint64, page uint64) ([]*models.Comment, error) {
	op := "CommentRepository.getPostComments"

	offset := (page - 1) * size

	query := r.qBuilder.
		Select(commentsAllCol).
		From(database.CommentsTable).
		Where(squirrel.Eq{commentsPostIdCol: postId}).
		Limit(size).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.GenerateSqlErr(err, op)
	}

	comments := make([]*models.Comment, 0)
	if err := r.db.SelectContext(ctx, &comments, sql, args...); err != nil {
		return nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return comments, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, commentId uuid.UUID) error {
	op := "CommentRepository.deleteComment"

	query := r.qBuilder.
		Delete(database.CommentsTable).
		Where(squirrel.Eq{commentsIdCol: commentId})

	sql, args, err := query.ToSql()
	if err != nil {
		return rep_utils.GenerateSqlErr(err, op)
	}

	if _, err := r.db.ExecContext(ctx, sql, args...); err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	return nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, updateData repositories_transfer.UpdateCommentInfo) error {
	op := "CommentRepository.updateComment"

	query := r.qBuilder.Update(database.CommentsTable).Where(squirrel.Eq{commentsIdCol: updateData.Id})

	//TODO
	for _, item := range updateData.UpdateData {
		if item.Name == commentsIdCol || item.Name == commentsUserIdCol || item.Name == commentsPostIdCol {
			continue
		}
		query = query.Set(item.Name, item.Value)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return rep_utils.GenerateSqlErr(err, op)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	return nil
}
