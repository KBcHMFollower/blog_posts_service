package repository

import (
	"context"
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/database"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PostRepository struct {
	db database.DBWrapper
}

func NewPostRepository(db database.DBWrapper) (*PostRepository, error) {
	return &PostRepository{
		db: db,
	}, nil
}

func (r *PostRepository) CreatePost(ctx context.Context, createData repositories_transfer.CreatePostInfo) (uuid.UUID, *models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	fmt.Println(createData)

	post := models.CreatePost(
		createData.UserId,
		createData.Title,
		createData.TextContent,
		createData.ImagesContent)

	fmt.Println(post)

	query := builder.
		Insert(POSTS_TABLE).
		Columns(ID_FIELD, USER_ID_FIELD, "title", "text_content", "images_content", "created_at").
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

	getSql, getArgs, err := builder.
		Select("*").
		From(POSTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: insertId}).
		ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	row := r.db.QueryRowContext(ctx, getSql, getArgs...)

	fmt.Println("1")

	var createdPost models.Post
	err = row.Scan(createdPost.GetPointersArray()...)

	fmt.Println(createdPost)
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	return createdPost.Id, &createdPost, nil
}

func (r *PostRepository) GetPost(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Select("*").
		From(POSTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: id})

	sql, args, _ := query.ToSql()

	var post models.Post

	row := r.db.QueryRowContext(ctx, sql, args...)
	err := row.Scan(post.GetPointersArray()...)
	if err != nil {
		return nil, fmt.Errorf("can`t scan properties from db : %v", err)
	}

	return &post, nil
}

func (r *PostRepository) GetPostsByUserId(ctx context.Context, getInfo repositories_transfer.GetPostByUserIdInfo) ([]*models.Post, uint, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	posts := make([]*models.Post, 0)

	offset := (getInfo.Page - 1) * getInfo.Size

	query := builder.
		Select("*").
		From(POSTS_TABLE).
		Where(squirrel.Eq{USER_ID_FIELD: getInfo.UserId}).
		Limit(uint64(getInfo.Size)).
		Offset(uint64(offset))

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

		err := rows.Scan(post.GetPointersArray()...)
		if err != nil {
			return posts, 0, fmt.Errorf("error in parse post from db: %v", err)
		}

		posts = append(posts, &post)
	}

	countQuery := builder.
		Select("COUNT(*)").
		From(POSTS_TABLE).
		Where(squirrel.Eq{USER_ID_FIELD: getInfo.UserId})

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

func (r *PostRepository) DeletePost(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Delete(POSTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	getSql, getArgs, err := builder.Select("*").
		From(POSTS_TABLE).
		Where(squirrel.Eq{ID_FIELD: id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	var post models.Post
	getRow := r.db.QueryRowContext(ctx, getSql, getArgs...)
	err = getRow.Scan(post.GetPointersArray()...)
	if err != nil {
		return nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error in execute sql-query : %v", err)
	}
	defer rows.Close()

	return &post, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, updateData repositories_transfer.UpdatePostInfo) (*models.Post, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Update(POSTS_TABLE).
		Where("id = ?", updateData.Id)

	for _, item := range updateData.UpdateData {
		if item.Name == ID_FIELD || item.Name == USER_ID_FIELD {
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

	queryGetPost := builder.
		Select("*").
		From(POSTS_TABLE).
		Where("id = ?", updateData.Id)
	sqlGetPost, argsGetPost, _ := queryGetPost.ToSql()

	row := r.db.QueryRowContext(ctx, sqlGetPost, argsGetPost...)

	var post models.Post
	err = row.Scan(post.GetPointersArray()...)
	if err != nil {
		return nil, fmt.Errorf("error scanning updated post : %v", err)
	}

	return &post, nil
}

func (r *PostRepository) DeleteUserPosts(ctx context.Context, userId uuid.UUID) error { //TODO
	op := "PostRepository.DeleteUserPosts"
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	query := builder.
		Delete(POSTS_TABLE).
		Where(squirrel.Eq{USER_ID_FIELD: userId})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error in generate sql-query : %v", err)
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error in execute sql-query : %v", err)
	}

	insertBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insertQuery := insertBuilder.Insert("transaction_events").
		SetMap(map[string]interface{}{
			"event_id":   uuid.New(),
			"event_type": "postsDeleted",
			"payload":    "",
		})

	sql, args, err = insertQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	return nil
}
