package app

import (
	"github.com/KBcHMFollower/test_plate_blog_service/config"
	database2 "github.com/KBcHMFollower/test_plate_blog_service/database"
	commentservice "github.com/KBcHMFollower/test_plate_blog_service/internal/services/comment_service"
	postService "github.com/KBcHMFollower/test_plate_blog_service/internal/services/post_service"
	"log/slog"

	grpcapp "github.com/KBcHMFollower/test_plate_blog_service/internal/app/grpc"
	"github.com/KBcHMFollower/test_plate_blog_service/internal/repository"
)

type App struct {
	GRPCServer *grpcapp.GRPCApp
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	op := "App.New"
	appLog := log.With(
		slog.String("op", op),
	)

	dbDriver, db, err := database2.New(cfg.Storage.ConnectionString)
	if err != nil {
		appLog.Error("db connection error: ", err)
		panic(err)
	}

	postRepository, err := repository.NewPostRepository(dbDriver)
	if err != nil {
		appLog.Error("TODO:", err)
		panic(err)
	}

	commRepository := repository.NewCommentRepository(dbDriver)

	if err := database2.ForceMigrate(db, cfg.Storage.MigrationPath); err != nil {
		appLog.Error("db migrate error: ", err)
		panic(err)
	}

	postService := postService.New(postRepository, log)
	commService := commentservice.New(commRepository, log)

	GRPCApp := grpcapp.New(cfg.GRpc.Port, log, postService, commService)

	return &App{
		GRPCServer: GRPCApp,
	}
}
