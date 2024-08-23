package main

import (
	"github.com/KBcHMFollower/blog_posts_service/internal/app"
	"github.com/KBcHMFollower/blog_posts_service/internal/config"
	sloglogger "github.com/KBcHMFollower/blog_posts_service/internal/logger/slog"
)

func main() {
	cfg := config.MustLoad()

	log := sloglogger.ConfigureLogger(cfg.Env)

	webApp := app.New(log, cfg)

	webApp.Run()

	log.Info("server started with: ")
}
