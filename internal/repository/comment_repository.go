package repository

import (
	"context"
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	POSTS_TABLE    = "posts"
	COMMENTS_TABLE = "comments"
	ID_FIELD       = "id"
	USER_ID_FIELD  = "user_id"
)

type CommentRepository struct {
	db database.DBWrapper
}

func NewCommentRepository(db database.DBWrapper) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, createData repositories_transfer.CreateCommentInfo) (uuid.UUID, *models.Comment, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	comment := models.CreateComment(createData.PostId, createData.UserId, createData.Content)

	query := builder.
		Insert(COMMENTS_TABLE).
		Columns(ID_FIELD, USER_ID_FIELD, "post_id", "content").
		Values(comment.Id, comment.UserId, comment.PostId, comment.Content).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {

		return uuid.Nil, nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	var insertId string

	idRow := r.db.QueryRowContext(ctx, sql, args...)

	if err := idRow.Scan(&insertId); err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	getSql, getArgs, err := builder.Select("*").
		From(COMMENTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: insertId}).
		ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	row := r.db.QueryRowContext(ctx, getSql, getArgs...)

	var createdComment models.Comment
	err = row.Scan(createdComment.GetPointerArray()...)
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	return createdComment.Id, &createdComment, nil
}

func (r *CommentRepository) GetComment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Select("*").
		From(COMMENTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: commentId})

	sql, args, _ := query.ToSql()

	var comment models.Comment

	row := r.db.QueryRowContext(ctx, sql, args...)
	err := row.Scan(comment.GetPointerArray()...)
	if err != nil {
		return nil, fmt.Errorf("can`t scan properties from db : %v", err)
	}

	return &comment, nil
}

func (r *CommentRepository) GetPostComments(ctx context.Context, postId uuid.UUID, size uint64, page uint64) ([]*models.Comment, uint, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	comments := make([]*models.Comment, 0)

	offset := (page - 1) * size

	query := builder.
		Select("*").
		From(COMMENTS_TABLE).
		Where(squirrel.Eq{"post_id": postId}).
		Limit(size).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return comments, 0, fmt.Errorf("error in generate sql-query : %v", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return comments, 0, fmt.Errorf("error in quey for db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment

		err := rows.Scan(comment.GetPointerArray()...)
		if err != nil {
			return comments, 0, fmt.Errorf("error in parse post from db: %v", err)
		}

		comments = append(comments, &comment)
	}

	countQuery := builder.
		Select("COUNT(*)").
		From(COMMENTS_TABLE).
		Where(squirrel.Eq{"post_id": postId})

	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return comments, 0, fmt.Errorf("error in generate sql-query : %v", err)
	}

	var totalCount uint

	countRow := r.db.QueryRowContext(ctx, countSql, countArgs...)
	if err := countRow.Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("can`t scan properties from db : %v", err)
	}

	return comments, totalCount, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, commentId uuid.UUID) (*models.Comment, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Delete(COMMENTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: commentId})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	getSql, getArgs, err := builder.Select("*").
		From(COMMENTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: commentId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	var comment models.Comment
	getRow := r.db.QueryRowContext(ctx, getSql, getArgs...)
	err = getRow.Scan(comment.GetPointerArray()...)
	if err != nil {
		return nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error in execute sql-query : %v", err)
	}
	defer rows.Close()

	return &comment, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, updateData repositories_transfer.UpdateCommentInfo) (*models.Comment, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.Update(COMMENTS_TABLE).Where(squirrel.Eq{ID_FIELD: updateData.Id})

	for _, item := range updateData.UpdateData {
		if item.Name == ID_FIELD || item.Name == USER_ID_FIELD || item.Name == "post_id" {
			continue
		}
		query = query.Set(item.Name, item.Value)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error in execute sql-query : %v", err)
	}

	queryGetComment := builder.Select("*").From(COMMENTS_TABLE).Where("id = ?", updateData.Id)
	sqlGetComment, argsGetComment, _ := queryGetComment.ToSql()

	row := r.db.QueryRowContext(ctx, sqlGetComment, argsGetComment...)

	var comment models.Comment
	err = row.Scan(comment.GetPointerArray()...)
	if err != nil {
		return nil, fmt.Errorf("error scanning updated post : %v", err)
	}

	return &comment, nil
}
