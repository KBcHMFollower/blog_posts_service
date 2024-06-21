package main

import (
	"github.com/KBcHMFollower/test_plate_blog_service/config"
	"github.com/KBcHMFollower/test_plate_blog_service/internal/app"
	sloglogger "github.com/KBcHMFollower/test_plate_blog_service/internal/logger/slog"
)

func main() {
	cfg := config.MustLoad()

	log := sloglogger.ConfigureLogger(cfg.Env)

	app := app.New(log, cfg)

	if err := app.GRPCServer.Run(); err != nil {
		panic(err)
	}

	log.Info("server started with: ")
}
