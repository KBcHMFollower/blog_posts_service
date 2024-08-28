package rep_utils

import (
	"fmt"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
)

func GetExecutor(r database.DBWrapper, tx database.Transaction) database.Executor { //TODO: ПОДУМАТЬ ОБ ЭТОМ, ПИЗДАТЕНЬКО ВЫШЛО
	if tx == nil {
		return r
	}
	return tx
}

func GenerateSqlErr(err error, op string) error {
	return fmt.Errorf("%s: error in generating sql: %w", op, err)
}

func ExecuteSqlErr(err error, op string) error {
	return fmt.Errorf("%s: error in executing sql: %w", op, err)
}
