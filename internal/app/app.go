package app

import (
	"log/slog"

	grpcapp "github.com/KBcHMFollower/test_plate_user_service/internal/app/grpc"
	"github.com/KBcHMFollower/test_plate_user_service/internal/config"
	"github.com/KBcHMFollower/test_plate_user_service/internal/database"
	"github.com/KBcHMFollower/test_plate_user_service/internal/repository"
	postService "github.com/KBcHMFollower/test_plate_user_service/internal/services"
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

	dbDriver, db, err := database.New(cfg.Storage.ConnectionString)
	if err != nil {
		appLog.Error("db connection error: ", err)
		panic(err)
	}

	postRepository, err := repository.NewPostRepository(dbDriver)
	if err != nil {
		appLog.Error("TODO:", err)
		panic(err)
	}
	if err := database.ForceMigrate(db, cfg.Storage.MigrationPath); err != nil {
		appLog.Error("db migrate error: ", err)
		panic(err)
	}

	postService := postService.New(postRepository, log)

	GRPCApp := grpcapp.New(cfg.GRpc.Port, log, postService)

	return &App{
		GRPCServer: GRPCApp,
	}
}
