package repository

import (
	"context"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	repositoriestransfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	reputils "github.com/KBcHMFollower/blog_posts_service/internal/repository/lib"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	postsTable = "posts"
)

const (
	postsIdCol            = "id"
	postsUserIdCol        = "user_id"
	postsTitleCol         = "title"
	postsTextContentCol   = "text_content"
	postsImagesContentCol = "images_content"
	postsCreatedAtCol     = "created_at"
	postsAllCol           = "*"
	postsSqlCount         = "COUNT(*)"
)

type PostRepository struct {
	db       database.DBWrapper
	qBuilder squirrel.StatementBuilderType //TODO: СЛИШКОМ ПРЯМАЯ ЗАВИСИМОСТЬ
}

func NewPostRepository(db database.DBWrapper) *PostRepository {
	return &PostRepository{
		db:       db,
		qBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostRepository) Create(ctx context.Context, createData repositoriestransfer.CreatePostInfo) (uuid.UUID, error) {
	post := models.CreatePost(
		createData.UserId,
		createData.Title,
		createData.TextContent,
		createData.ImagesContent)

	query := r.qBuilder.
		Insert(postsTable).
		SetMap(map[string]interface{}{
			postsIdCol:            post.Id,
			postsUserIdCol:        post.UserId,
			postsTitleCol:         post.Title,
			postsTextContentCol:   post.TextContent,
			postsImagesContentCol: post.ImagesContent,
		}).
		Suffix("RETURNING \"id\"")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return uuid.New(), reputils.ReturnGenerateSqlError(ctx, err)
	}

	var insertId uuid.UUID
	if err := r.db.GetContext(ctx, &insertId, sqlStr, args...); err != nil {
		return uuid.New(), reputils.ReturnExecuteSqlError(ctx, err)
	}

	return insertId, nil
}

func (r *PostRepository) Post(ctx context.Context, info repositoriestransfer.GetPostInfo) (*models.Post, error) {
	query := r.qBuilder.
		Select(postsAllCol).
		From(postsTable).
		Where(squirrel.Eq(reputils.ConvertMapKeysToStrings(info.Condition)))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, reputils.ReturnGenerateSqlError(ctx, err)
	}

	var post models.Post

	if err := r.db.GetContext(ctx, &post, sql, args...); err != nil {
		return nil, reputils.ReturnExecuteSqlError(ctx, err)
	}

	return &post, nil
}

func (r *PostRepository) Count(ctx context.Context, info repositoriestransfer.GetPostsCountInfo) (uint64, error) {
	countQuery := r.qBuilder.
		Select(postsSqlCount).
		From(postsTable).
		Where(squirrel.Eq(reputils.ConvertMapKeysToStrings(info.Condition)))

	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return 0, reputils.ReturnGenerateSqlError(ctx, err)
	}

	var totalCount uint64

	if err := r.db.GetContext(ctx, &totalCount, countSql, countArgs...); err != nil {
		return 0, reputils.ReturnExecuteSqlError(ctx, err)
	}

	return totalCount, nil
}

func (r *PostRepository) Posts(ctx context.Context, getInfo repositoriestransfer.GetPostsInfo) ([]*models.Post, error) {
	offset := (getInfo.Page - 1) * getInfo.Size

	query := r.qBuilder.
		Select(postsAllCol).
		From(postsTable).
		Where(squirrel.Eq(reputils.ConvertMapKeysToStrings(getInfo.Condition))).
		Limit(uint64(getInfo.Size)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, reputils.ReturnGenerateSqlError(ctx, err)
	}

	posts := make([]*models.Post, 0)
	if err := r.db.SelectContext(ctx, &posts, sqlStr, args...); err != nil {
		return nil, reputils.ReturnExecuteSqlError(ctx, err)
	}

	return posts, nil
}

func (r *PostRepository) Delete(ctx context.Context, info repositoriestransfer.DeletePostsInfo, tx database.Transaction) error {
	executor := reputils.GetExecutor(r.db, tx)

	query := r.qBuilder.
		Delete(postsTable).
		Where(squirrel.Eq(reputils.ConvertMapKeysToStrings(info.Condition)))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return reputils.ReturnGenerateSqlError(ctx, err)
	}

	if _, err := executor.ExecContext(ctx, sqlStr, args...); err != nil {
		return reputils.ReturnExecuteSqlError(ctx, err)
	}

	return nil
}

func (r *PostRepository) Update(ctx context.Context, updateData repositoriestransfer.UpdatePostInfo) error {
	query := r.qBuilder.
		Update(postsTable).
		Where(squirrel.Eq(reputils.ConvertMapKeysToStrings(updateData.Condition))).
		SetMap(reputils.ConvertMapKeysToStrings(updateData.UpdateData))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return reputils.ReturnGenerateSqlError(ctx, err)
	}

	if _, err = r.db.ExecContext(ctx, sqlStr, args...); err != nil {
		return reputils.ReturnExecuteSqlError(ctx, err)
	}

	return nil
}
