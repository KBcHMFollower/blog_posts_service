package postgresrepository

import (
	"context"
	"fmt"

	"github.com/KBcHMFollower/test_plate_user_service/internal/domain/models"
	"github.com/KBcHMFollower/test_plate_user_service/internal/repository"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *PostgresRepository) CreatePost(ctx context.Context, createData repository.CreatePostData) (uuid.UUID, *models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	post := models.CreatePost(
		createData.User_id,
		createData.Title,
		createData.TextContent,
		*createData.ImagesContent)

	query := builder.
		Insert("posts").
		Columns("id", "user_id", "title", "text_content", "images_content", "created_at").
		Values(post.Id, post.User_id, post.Title, post.TextContent, post.ImageContent, post.Created_at)

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	row := r.db.QueryRowContext(ctx, sql, args...)

	var createdPost models.Post
	err = row.Scan(
		&createdPost.Id,
		&createdPost.User_id,
		&createdPost.Title,
		&createdPost.TextContent,
		&createdPost.ImageContent,
		&createdPost.Created_at)
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	return createdPost.Id, &createdPost, nil
}

func (r *PostgresRepository) GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Select("*").
		From("posts").
		Where(squirrel.Eq{"id": id})

	sql, args, _ := query.ToSql()

	var post models.Post

	row := r.db.QueryRowContext(ctx, sql, args...)
	row.Scan(&post.Id,
		&post.User_id,
		&post.Title,
		&post.TextContent,
		&post.ImageContent, &post.Created_at)

	return &post, nil
}

func (r *PostgresRepository) GetPostsByUserId(ctx context.Context, user_id int) ([]*models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Select("*").
		From("posts").
		Where(squirrel.Eq{"user_id": user_id})

	sql, args, _ := query.ToSql()

	posts := make([]*models.Post, 0)

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return posts, fmt.Errorf("error in quey for db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.Id,
			&post.User_id,
			&post.Title,
			&post.TextContent,
			&post.ImageContent, &post.Created_at)
		if err != nil {
			return posts, fmt.Errorf("error in parse post from db: %v", err)
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostgresRepository) DeletePost(ctx context.Context, id uuid.UUID) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Delete("posts").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error in generate sql-query : %v", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	defer rows.Close()

	if err != nil {
		return fmt.Errorf("error in execute sql-query : %v", err)
	}

	return nil
}

func (r *PostgresRepository) UpdatePost(ctx context.Context, updateData repository.UpdatePostData) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.Update("posts").Where("id = $", updateData.Id)

	for _, item := range updateData.UpdateData {
		query = query.Set(item.Name, item.Value)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error in generate sql-query : %v", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args)
	defer rows.Close()

	if err != nil {
		return fmt.Errorf("error in execute sql-query : %v", err)
	}

	return nil
}
