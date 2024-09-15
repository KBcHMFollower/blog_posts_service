package services_dep

import (
	"context"
	"database/sql"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
)

type TransactionCreator interface {
	BeginTxCtx(ctx context.Context, opts *sql.TxOptions) (database.Transaction, error)
}
