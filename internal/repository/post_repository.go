package repository

import (
	"context"
	"fmt"

	"github.com/KBcHMFollower/test_plate_user_service/internal/database"
	"github.com/KBcHMFollower/test_plate_user_service/internal/domain/models"
	"github.com/KBcHMFollower/test_plate_user_service/internal/lib/consts"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostRepository struct {
	db database.DBWrapper
}

func NewPostRepository(db database.DBWrapper) (*PostRepository, error) {
	return &PostRepository{
		db: db,
	}, nil
}

func (r *PostRepository) CreatePost(ctx context.Context, createData CreatePostData) (uuid.UUID, *models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	post := models.CreatePost(
		createData.User_id,
		createData.Title,
		createData.TextContent,
		*createData.ImagesContent)

	query := builder.
		Insert(consts.POSTS_TABLE).
		Columns(consts.ID_FIELD, consts.USER_ID_FIELD, "title", "text_content", "images_content", "created_at").
		Values(post.Id, post.UserId, post.Title, post.TextContent, post.ImagesContent, post.CreatedAt).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	var insertId string

	idRow := r.db.QueryRowContext(ctx, sql, args...)

	if err := idRow.Scan(&insertId); err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	getSql, getArgs, err := builder.Select("*").From(consts.POSTS_TABLE).Where(squirrel.Eq{consts.ID_FIELD: insertId}).ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	row := r.db.QueryRowContext(ctx, getSql, getArgs...)

	var createdPost models.Post
	err = row.Scan(
		&createdPost.Id,
		&createdPost.UserId,
		&createdPost.Title,
		&createdPost.TextContent,
		&createdPost.ImagesContent,
		&createdPost.CreatedAt)
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	return createdPost.Id, &createdPost, nil
}

func (r *PostRepository) GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Select("*").
		From(consts.POSTS_TABLE).
		Where(squirrel.Eq{consts.ID_FIELD: id})

	sql, args, _ := query.ToSql()

	var post models.Post

	row := r.db.QueryRowContext(ctx, sql, args...)
	err := row.Scan(&post.Id,
		&post.UserId,
		&post.Title,
		&post.TextContent,
		&post.ImagesContent,
		&post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("can`t scan properties from db : %v", err)
	}

	return &post, nil
}

func (r *PostRepository) GetPostsByUserId(ctx context.Context, user_id uuid.UUID, size uint64, page uint64) ([]*models.Post, uint, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	posts := make([]*models.Post, 0)

	offset := (page - 1) * size

	query := builder.
		Select("*").
		From(consts.POSTS_TABLE).
		Where(squirrel.Eq{consts.USER_ID_FIELD: user_id}).
		Limit(size).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return posts, 0, fmt.Errorf("error in generate sql-query : %v", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return posts, 0, fmt.Errorf("error in quey for db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.Id,
			&post.UserId,
			&post.Title,
			&post.TextContent,
			&post.ImagesContent, &post.CreatedAt)
		if err != nil {
			return posts, 0, fmt.Errorf("error in parse post from db: %v", err)
		}

		posts = append(posts, &post)
	}

	countQuery := builder.
		Select("COUNT(*)").
		From(consts.POSTS_TABLE)

	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return posts, 0, fmt.Errorf("error in generate sql-query : %v", err)
	}

	var totalCount uint

	countRow := r.db.QueryRowContext(ctx, countSql, countArgs...)
	if err := countRow.Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("can`t scan properties from db : %v", err)
	}

	return posts, totalCount, nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id uuid.UUID) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Delete(consts.POSTS_TABLE).
		Where(squirrel.Eq{consts.ID_FIELD: id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error in generate sql-query : %v", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error in execute sql-query : %v", err)
	}
	defer rows.Close()

	return nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, updateData UpdatePostData) (*models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.Update(consts.POSTS_TABLE).Where("id = ?", updateData.Id)

	for _, item := range updateData.UpdateData {
		if item.Name == consts.ID_FIELD || item.Name == consts.USER_ID_FIELD {
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

	queryGetPost := builder.Select("*").From(consts.POSTS_TABLE).Where("id = ?", updateData.Id)
	sqlGetPost, argsGetPost, _ := queryGetPost.ToSql()

	row := r.db.QueryRowContext(ctx, sqlGetPost, argsGetPost...)

	var post models.Post
	err = row.Scan(&post.Id, &post.UserId, &post.Title, &post.TextContent, &post.ImagesContent, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning updated post : %v", err)
	}

	return &post, nil
}
