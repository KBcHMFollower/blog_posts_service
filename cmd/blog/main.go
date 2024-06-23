package main

import (
	"github.com/KBcHMFollower/blog_posts_service/config"
	"github.com/KBcHMFollower/blog_posts_service/internal/app"
	sloglogger "github.com/KBcHMFollower/blog_posts_service/internal/logger/slog"
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
