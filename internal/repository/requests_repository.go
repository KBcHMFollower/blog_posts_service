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

type RequestsStore interface {
	Create(ctx context.Context, key uuid.UUID, payload string) (uuid.UUID, *models.Request, error)
	Get(ctx context.Context, key uuid.UUID) (*models.Request, error)
}

type RequestsRepository struct {
	db database.DBWrapper
}

func NewRequestsRepository(db database.DBWrapper) (*RequestsRepository, error) {
	return &RequestsRepository{
		db: db,
	}, nil
}

func (r *RequestsRepository) Create(ctx context.Context, info repositories_transfer.CreateRequestInfo) (uuid.UUID, *models.Request, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	request := models.Request{
		Id:             uuid.New(),
		IdempotencyKey: info.Key,
		Payload:        info.Payload,
	}

	query := builder.
		Insert("requests").
		Columns(ID_FIELD, "idempotency_key", "payload").
		Values(request.Id, request.IdempotencyKey, request.Payload).
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
		From("requests").
		Where(squirrel.Eq{ID_FIELD: insertId}).
		ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	row := r.db.QueryRowContext(ctx, getSql, getArgs...)

	fmt.Println("1")

	var createdRequest models.Request
	err = row.Scan(createdRequest.GetPointersArray()...)

	fmt.Println(createdRequest)
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	return createdRequest.Id, &createdRequest, nil
}

func (r *PostRepository) Get(ctx context.Context, key uuid.UUID) (*models.Request, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := builder.
		Select("*").
		From("requests").
		Where(squirrel.Eq{"idempotency_key": key})

	sql, args, _ := query.ToSql()

	var request models.Request

	row := r.db.QueryRowContext(ctx, sql, args...)
	err := row.Scan(request.GetPointersArray()...)
	if err != nil {
		return nil, fmt.Errorf("can`t scan properties from db : %v", err)
	}

	return &request, nil
}
