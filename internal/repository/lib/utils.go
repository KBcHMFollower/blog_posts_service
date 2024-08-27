package rep_utils

import (
	"github.com/KBcHMFollower/blog_posts_service/internal/database"
)

func GetExecutor(r database.DBWrapper, tx database.Transaction) database.Executor { //TODO: ПОДУМАТЬ ОБ ЭТОМ, ПИЗДАТЕНЬКО ВЫШЛО
	if tx == nil {
		return r
	}
	return tx
}
