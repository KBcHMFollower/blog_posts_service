package repository

import (
	"context"
	"database/sql"
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
	rKeysAllCol            = "*"
)

type RequestsStore interface {
	Create(ctx context.Context, info repositories_transfer.CreateRequestInfo, tx *sql.Tx) (uuid.UUID, *models.Request, error)
	Get(ctx context.Context, key uuid.UUID, tx *sql.Tx) (*models.Request, error)
}

type RequestsRepository struct {
	db       database.DBWrapper
	qBuilder squirrel.StatementBuilderType //TODO: СЛИШКОМ ПРЯМАЯ ЗАВИСИМОСТЬ
}

func NewRequestsRepository(db database.DBWrapper) *RequestsRepository {
	return &RequestsRepository{
		db:       db,
		qBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *RequestsRepository) Create(ctx context.Context, info repositories_transfer.CreateRequestInfo, tx database.Transaction) (uuid.UUID, error) {
	op := "RequestRepository.create"

	executor := rep_utils.GetExecutor(r.db, tx)

	request := models.Request{
		Id:             uuid.New(),
		IdempotencyKey: info.Key,
	}

	query := r.qBuilder.
		Insert(database.RequestKeysTable).
		SetMap(map[string]interface{}{
			rKeysIdCol:             request.Id,
			rKeysIdempotencyKeyCol: request.IdempotencyKey,
		}).
		Suffix("RETURNING \"id\"")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return uuid.New(), rep_utils.GenerateSqlErr(err, op)
	}

	var insertId uuid.UUID
	if err := executor.GetContext(ctx, &insertId, sqlStr, args...); err != nil {
		return uuid.New(), rep_utils.ExecuteSqlErr(err, op)
	}

	return insertId, nil
}

func (r *RequestsRepository) Get(ctx context.Context, key uuid.UUID, tx database.Transaction) (*models.Request, error) {
	op := "RequestRepository.get"

	executor := rep_utils.GetExecutor(r.db, tx)

	query := r.qBuilder.
		Select(rKeysAllCol).
		From(database.RequestKeysTable).
		Where(squirrel.Eq{rKeysIdempotencyKeyCol: key})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, rep_utils.GenerateSqlErr(err, op)
	}

	var request models.Request

	if err := executor.GetContext(ctx, &request, sqlStr, args...); err != nil {
		return nil, rep_utils.ExecuteSqlErr(err, op)
	}

	return &request, nil
}
