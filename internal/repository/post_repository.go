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

// TODO: ДОВЫНОСИТЬ ВСЕ В КОНСТАНТЫ
// TODO: ВРЯД ЛИ НОРМ БРАТЬ СРАЗУ ВСЕ ПОЛЯ С ТАБЛИЦЫ
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

func (r *PostRepository) CreatePost(ctx context.Context, createData repositories_transfer.CreatePostInfo) (uuid.UUID, error) {
	op := "PostRepository.createPost"

	post := models.CreatePost(
		createData.UserId,
		createData.Title,
		createData.TextContent,
		createData.ImagesContent)

	query := r.qBuilder.
		Insert(database.PostsTable).
		SetMap(map[string]interface{}{
			postsIdCol:            post.Id,
			postsUserIdCol:        post.UserId,
			postsTitleCol:         post.Title,
			postsTextContentCol:   post.TextContent,
			postsImagesContentCol: post.ImagesContent,
			postsCreatedAtCol:     post.CreatedAt,
		}).
		Suffix("RETURNING \"id\"")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return uuid.New(), rep_utils.GenerateSqlErr(err, op)
	}

	var insertId uuid.UUID
	if err := r.db.GetContext(ctx, &insertId, sqlStr, args...); err != nil {
		return uuid.New(), rep_utils.ExecuteSqlErr(err, op)
	}

	return insertId, nil
}

func (r *PostRepository) GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	op := "PostRepository.getPost"

	query := r.qBuilder.
		Select(postsAllCol).
		From(database.PostsTable).
		Where(squirrel.Eq{postsIdCol: id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.GenerateSqlErr(err, op)
	}

	var post models.Post

	if err := r.db.GetContext(ctx, &post, sql, args...); err != nil {
		return nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return &post, nil
}

func (r *PostRepository) GetUserPostsCount(ctx context.Context, userId uuid.UUID) (uint, error) {
	op := "PostRepository.getUserPostsCount"

	countQuery := r.qBuilder.
		Select(postsSqlCount).
		From(database.PostsTable).
		Where(squirrel.Eq{postsUserIdCol: userId})

	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return 0, rep_utils.GenerateSqlErr(err, op)
	}

	var totalCount uint

	if err := r.db.GetContext(ctx, &totalCount, countSql, countArgs...); err != nil {
		return 0, rep_utils.ExecuteSqlErr(err, op)
	}

	return totalCount, nil
}

func (r *PostRepository) GetPostsByUserId(ctx context.Context, getInfo repositories_transfer.GetPostByUserIdInfo) ([]*models.Post, error) {
	op := "PostRepository.getPostsByUserId"

	offset := (getInfo.Page - 1) * getInfo.Size

	query := r.qBuilder.
		Select(postsAllCol).
		From(database.PostsTable).
		Where(squirrel.Eq{postsUserIdCol: getInfo.UserId}).
		Limit(uint64(getInfo.Size)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.GenerateSqlErr(err, op)
	}

	posts := make([]*models.Post, 0)
	if err := r.db.SelectContext(ctx, &posts, sqlStr, args...); err != nil {
		return nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return posts, nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id uuid.UUID) error {
	op := "PostRepository.deletePost"

	query := r.qBuilder.
		Delete(database.PostsTable).
		Where(squirrel.Eq{postsIdCol: id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return rep_utils.GenerateSqlErr(err, op)
	}

	if _, err := r.db.ExecContext(ctx, sqlStr, args...); err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	return nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, updateData repositories_transfer.UpdatePostInfo) error {
	op := "PostRepository.updatePost"

	query := r.qBuilder.
		Update(database.PostsTable).
		Where("id = ?", updateData.Id) //TODO

	for _, item := range updateData.UpdateData {
		if item.Name == postsIdCol || item.Name == postsUserIdCol {
			continue
		}
		query = query.Set(item.Name, item.Value)
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return rep_utils.GenerateSqlErr(err, op)
	}

	if _, err = r.db.ExecContext(ctx, sqlStr, args...); err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	return nil
}

func (r *PostRepository) DeleteUserPosts(ctx context.Context, userId uuid.UUID, tx database.Transaction) error {
	op := "PostRepository.deleteUserPosts"

	executor := rep_utils.GetExecutor(r.db, tx)

	query := r.qBuilder.
		Delete(database.PostsTable).
		Where(squirrel.Eq{postsUserIdCol: userId})

	toSql, args, err := query.ToSql()
	if err != nil {
		return rep_utils.GenerateSqlErr(err, op)
	}

	if _, err = executor.ExecContext(ctx, toSql, args...); err != nil {
		return rep_utils.ExecuteSqlErr(err, op)
	}

	return nil
}
