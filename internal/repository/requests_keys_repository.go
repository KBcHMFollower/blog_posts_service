package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
	repositories_transfer "github.com/KBcHMFollower/blog_posts_service/internal/domain/layers_TOs/repositories"
	"github.com/KBcHMFollower/blog_posts_service/internal/domain/models"
	rep_utils "github.com/KBcHMFollower/blog_posts_service/internal/repository/lib"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	rKeysIdCol             = "id"
	rKeysIdempotencyKeyCol = "idempotency_key"
	rKeysPayloadCol        = "payload"
	rKeysAllCol            = "*"
)

type RequestsStore interface {
	Create(ctx context.Context, info repositories_transfer.CreateRequestInfo, tx *sql.Tx) (uuid.UUID, *models.Request, error)
	Get(ctx context.Context, key uuid.UUID, tx *sql.Tx) (*models.Request, error)
}

type RequestsRepository struct {
	db database.DBWrapper
}

func NewRequestsRepository(db database.DBWrapper) *RequestsRepository {
	return &RequestsRepository{
		db: db,
	}
}

func (r *RequestsRepository) Create(ctx context.Context, info repositories_transfer.CreateRequestInfo, tx *sql.Tx) (uuid.UUID, *models.Request, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	executor := rep_utils.GetExecutor(r.db, tx)

	request := models.Request{
		Id:             uuid.New(),
		IdempotencyKey: info.Key,
	}

	query := builder.
		Insert(database.RequestKeysTable).
		Columns(rKeysIdCol, rKeysIdempotencyKeyCol, rKeysPayloadCol).
		Values(request.Id, request.IdempotencyKey).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	var insertId string

	idRow := executor.QueryRowContext(ctx, sql, args...)

	if err := idRow.Scan(&insertId); err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	getSql, getArgs, err := builder.
		Select(rKeysAllCol).
		From(database.RequestKeysTable).
		Where(squirrel.Eq{rKeysIdCol: insertId}).
		ToSql()
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in generate sql-query : %v", err)
	}

	row := executor.QueryRowContext(ctx, getSql, getArgs...)

	fmt.Println("1")

	var createdRequest models.Request
	err = row.Scan(createdRequest.GetPointersArray()...)

	fmt.Println(createdRequest)
	if err != nil {
		return uuid.New(), nil, fmt.Errorf("error in scan property from db : %v", err)
	}

	return createdRequest.Id, &createdRequest, nil
}

func (r *RequestsRepository) Get(ctx context.Context, key uuid.UUID, tx *sql.Tx) (*models.Request, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	executor := rep_utils.GetExecutor(r.db, tx)

	query := builder.
		Select(rKeysAllCol).
		From(database.RequestKeysTable).
		Where(squirrel.Eq{rKeysIdempotencyKeyCol: key})

	sql, args, _ := query.ToSql()

	var request models.Request

	row := executor.QueryRowContext(ctx, sql, args...)
	err := row.Scan(request.GetPointersArray()...)
	if err != nil {
		return nil, fmt.Errorf("can`t scan properties from db : %v", err)
	}

	return &request, nil
}
