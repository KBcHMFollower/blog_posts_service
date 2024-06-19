package main

import (
	"github.com/KBcHMFollower/test_plate_user_service/internal/config"
	sloglogger "github.com/KBcHMFollower/test_plate_user_service/internal/logger/slog"
)

func main() {
	cfg := config.MustLoad()

	log := sloglogger.ConfigureLogger(cfg.Env)

	log.Info("server started with: ")
}
