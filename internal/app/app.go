package app

import (
	"fmt"
	amqpapp "github.com/KBcHMFollower/blog_posts_service/internal/app/amqp_app"
	"github.com/KBcHMFollower/blog_posts_service/internal/app/store_app"
	"github.com/KBcHMFollower/blog_posts_service/internal/app/workers_app"
	"github.com/KBcHMFollower/blog_posts_service/internal/config"
	ctxerrors "github.com/KBcHMFollower/blog_posts_service/internal/domain/errors"
	"github.com/KBcHMFollower/blog_posts_service/internal/interceptors"
	"github.com/KBcHMFollower/blog_posts_service/internal/lib/validators"
	"github.com/KBcHMFollower/blog_posts_service/internal/logger"
	commentservice "github.com/KBcHMFollower/blog_posts_service/internal/services"
	"github.com/KBcHMFollower/blog_posts_service/internal/services/lib/circuid_breaker"
	"github.com/KBcHMFollower/blog_posts_service/internal/workers"
	"google.golang.org/grpc"
	"time"

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
	log logger.Logger,
	cfg *config.Config,
) *App {
	storeApp, err := store_app.New(cfg.Storage)
	ContinueOrPanic(err)
	amqpApp, err := amqpapp.NewAmqpApp(cfg.RabbitMq, log)
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
	reqService := commentservice.NewRequestsService(reqRepository, log)

	workersApp.AddWorker(workers.NewEventChecker(
		amqpApp.Client,
		eventRepository,
		log,
		storeApp.PostgresStore.Store,
	))

	vldor, err := validators.NewValidator()
	ContinueOrPanic(err)

	circuitBreaker := circuid_breaker.NewCircuitBreaker().Configure(func(options *circuid_breaker.CBOptions) {
		options.IgnorableErrors = []error{
			ctxerrors.ErrNotFound,
			ctxerrors.ErrUnauthorized,
			ctxerrors.ErrConflict,
			ctxerrors.ErrBadRequest,
		}
		options.OpenConditions = circuid_breaker.OpenCondition{
			FailuresRate: 40,
			TimeInterval: time.Duration(100),
		}
		options.CloseConditions = circuid_breaker.CloseCondition{
			SuccessRate: 80,
			Duration:    time.Duration(100),
		}
	})
	circuitBreaker.OnChangeStateHook = func(from circuid_breaker.BreakerState, to circuid_breaker.BreakerState) {
		log.Info("circuit breaker state changed", "from", from, "to", to)
	}

	interceptorsChain := grpc.ChainUnaryInterceptor(
		interceptors.CircuitBreakerInterceptor(circuitBreaker),
		interceptors.ErrorHandlerInterceptor(),
		interceptors.ReqLoggingInterceptor(log),
		interceptors.IdempotencyInterceptor(reqService),
	)

	GRPCApp := grpcapp.New(cfg.GRpc.Port, log, postService, commService, vldor, interceptorsChain)

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
