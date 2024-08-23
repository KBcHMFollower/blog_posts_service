package app

import (
	"fmt"
	amqpapp "github.com/KBcHMFollower/blog_posts_service/internal/app/amqp_app"
	"github.com/KBcHMFollower/blog_posts_service/internal/app/store_app"
	"github.com/KBcHMFollower/blog_posts_service/internal/app/workers_app"
	"github.com/KBcHMFollower/blog_posts_service/internal/config"
	commentservice "github.com/KBcHMFollower/blog_posts_service/internal/services"
	"github.com/KBcHMFollower/blog_posts_service/internal/workers"
	"log/slog"

	grpcapp "github.com/KBcHMFollower/blog_posts_service/internal/app/grpc_app"
	"github.com/KBcHMFollower/blog_posts_service/internal/repository"
)

type App struct {
	gRPCApp    *grpcapp.GRPCApp
	storeApp   *store_app.StoreApp
	amqpApp    *amqpapp.AmqpApp
	workersApp *workers_app.WorkersApp
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	//op := "App.NewCommentService"
	//appLog := log.With(
	//	slog.String("op", op),
	//)

	storeApp, err := store_app.New(cfg.Storage)
	ContinueOrPanic(err)
	amqpApp, err := amqpapp.New(cfg.RabbitMq)
	ContinueOrPanic(err)
	workersApp := workers_app.New()

	postRepository := repository.NewPostRepository(storeApp.PostgresStore.Store)
	commRepository := repository.NewCommentRepository(storeApp.PostgresStore.Store)
	eventRepository := repository.NewEventRepository(storeApp.PostgresStore.Store)
	reqRepository := repository.NewRequestsRepository(storeApp.PostgresStore.Store)

	postService := commentservice.NewPostsService(
		postRepository,
		reqRepository,
		eventRepository,
		storeApp.PostgresStore.Store,
		log,
	)
	commService := commentservice.NewCommentService(commRepository, log)

	workersApp.AddWorker(workers.NewEventChecker(amqpApp.Client, eventRepository, log))

	GRPCApp := grpcapp.New(cfg.GRpc.Port, log, postService, commService)

	return &App{
		gRPCApp:    GRPCApp,
		storeApp:   storeApp,
		amqpApp:    amqpApp,
		workersApp: workersApp,
	}
}

func (app *App) Run() {
	err := app.storeApp.Run()
	ContinueOrPanic(err)

	err = app.gRPCApp.Run()
	ContinueOrPanic(err)

	err = app.amqpApp.Start()
	ContinueOrPanic(err)

	err = app.workersApp.Run()
	ContinueOrPanic(err)
}

func (app *App) Stop() error {
	if err := app.storeApp.Stop(); err != nil {
		return fmt.Errorf("stop store_app app: %w", err)
	}

	if err := app.amqpApp.Stop(); err != nil {
		return fmt.Errorf("stop amqp_app app: %w", err)
	}

	app.gRPCApp.Stop()
	app.workersApp.Stop()

	return nil
}

func ContinueOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}
