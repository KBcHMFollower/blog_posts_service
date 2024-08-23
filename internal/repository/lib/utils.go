package rep_utils

import (
	"database/sql"
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
)

func GetExecutor(r database.DBWrapper, tx *sql.Tx) database.Executor { //TODO: ПОДУМАТЬ ОБ ЭТОМ, ПИЗДАТЕНЬКО ВЫШЛО
	if tx == nil {
		return r
	}
	return tx
}
