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
	commentsTableName = "comments"
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

func (r *CommentRepository) Create(ctx context.Context, createData repositories_transfer.CreateCommentInfo) (uuid.UUID, error) {
	comment := models.CreateComment(createData.PostId, createData.UserId, createData.Content)

	query := r.qBuilder.
		Insert(commentsTableName).
		SetMap(map[string]interface{}{
			commentsIdCol:      comment.Id,
			commentsUserIdCol:  comment.UserId,
			commentsPostIdCol:  comment.PostId,
			commentsContentCol: comment.Content,
		}).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, rep_utils.ReturnGenerateSqlError(ctx, err)
	}

	var insertId uuid.UUID
	if err := r.db.GetContext(ctx, &insertId, sql, args...); err != nil {
		return uuid.Nil, rep_utils.ReturnExecuteSqlError(ctx, err)
	}

	return insertId, nil
}

func (r *CommentRepository) Comment(ctx context.Context, info repositories_transfer.GetCommentInfo) (*models.Comment, error) {
	query := r.qBuilder.
		Select(commentsAllCol).
		From(commentsTableName).
		Where(squirrel.Eq(rep_utils.ConvertMapKeysToStrings(info.Condition)))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.ReturnGenerateSqlError(ctx, err)
	}

	var comment models.Comment
	if err := r.db.GetContext(ctx, &comment, sqlStr, args...); err != nil {
		return nil, rep_utils.ReturnExecuteSqlError(ctx, err)
	}

	return &comment, nil
}

func (r *CommentRepository) Count(ctx context.Context, info repositories_transfer.GetCommentsCountInfo) (uint64, error) {
	query := r.qBuilder.
		Select(commentsSqlCount).
		From(commentsTableName).
		Where(squirrel.Eq(rep_utils.ConvertMapKeysToStrings(info.Condition)))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, rep_utils.ReturnGenerateSqlError(ctx, err)
	}

	var count uint64
	if err := r.db.GetContext(ctx, &count, sqlStr, args...); err != nil {
		return 0, rep_utils.ReturnExecuteSqlError(ctx, err)
	}

	return count, nil
}

func (r *CommentRepository) Comments(ctx context.Context, info repositories_transfer.GetCommentsInfo) ([]*models.Comment, error) {
	offset := (info.Page - 1) * info.Size

	query := r.qBuilder.
		Select(commentsAllCol).
		From(commentsTableName).
		Where(squirrel.Eq(rep_utils.ConvertMapKeysToStrings(info.Condition))).
		Limit(info.Size).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.ReturnGenerateSqlError(ctx, err)
	}

	comments := make([]*models.Comment, 0)
	if err := r.db.SelectContext(ctx, &comments, sql, args...); err != nil {
		return nil, rep_utils.ReturnExecuteSqlError(ctx, err)
	}

	return comments, nil
}

func (r *CommentRepository) Delete(ctx context.Context, info repositories_transfer.DeleteCommentInfo) error {
	query := r.qBuilder.
		Delete(commentsTableName).
		Where(squirrel.Eq(rep_utils.ConvertMapKeysToStrings(info.Condition)))

	sql, args, err := query.ToSql()
	if err != nil {
		return rep_utils.ReturnGenerateSqlError(ctx, err)
	}

	if _, err := r.db.ExecContext(ctx, sql, args...); err != nil {
		return rep_utils.ReturnExecuteSqlError(ctx, err)
	}

	return nil
}

func (r *CommentRepository) Update(ctx context.Context, updateData repositories_transfer.UpdateCommentInfo) error {
	query := r.qBuilder.
		Update(commentsTableName).
		Where(squirrel.Eq(rep_utils.ConvertMapKeysToStrings(updateData.Condition))).
		SetMap(rep_utils.ConvertMapKeysToStrings(updateData.UpdateData))

	sql, args, err := query.ToSql()
	if err != nil {
		return rep_utils.ReturnGenerateSqlError(ctx, err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return rep_utils.ReturnExecuteSqlError(ctx, err)
	}

	return nil
}
