package app

import (
	"log/slog"

	grpcapp "github.com/KBcHMFollower/test_plate_user_service/internal/app/grpc"
	"github.com/KBcHMFollower/test_plate_user_service/internal/config"
	postgresrepository "github.com/KBcHMFollower/test_plate_user_service/internal/repository/postgres"
	postService "github.com/KBcHMFollower/test_plate_user_service/internal/servicces"
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

	rep, err := postgresrepository.New(cfg.Storage.ConnectionString)
	if err != nil {
		appLog.Error("db connection error: ", err)
		panic(err)
	}

	postService := postService.New(rep, log)

	GRPCApp := grpcapp.New(cfg.GRpc.Port, log, postService)

	return &App{
		GRPCServer: GRPCApp,
	}
}
