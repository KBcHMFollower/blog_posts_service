package grpcapp

import (
	"fmt"
	handlers_dep "github.com/KBcHMFollower/blog_posts_service/internal/handlers/dep"
	grpcserver2 "github.com/KBcHMFollower/blog_posts_service/internal/handlers/grpc"
	commentservice "github.com/KBcHMFollower/blog_posts_service/internal/services"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type GRPCApp struct {
	log        *slog.Logger
	port       int
	GRPCServer *grpc.Server
}

func New(
	port int,
	log *slog.Logger,
	postService *commentservice.PostService,
	commService *commentservice.CommentsService,
	validator handlers_dep.Validator,
	interceptor grpc.ServerOption,
) *GRPCApp {
	gRpcServer := grpc.NewServer(interceptor)

	grpcserver2.RegisterPostServer(gRpcServer, postService)
	grpcserver2.RegisterCommentsServer(gRpcServer, commService)

	return &GRPCApp{
		log:        log,
		port:       port,
		GRPCServer: gRpcServer,
	}
}

func (g *GRPCApp) Run() error {
	const op = "GRPCApp.Run"
	log := g.log.With(
		slog.String("op", op),
	)

	log.Info("grpc_app server is trying to get up")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		log.Error("server startup error ", err)
		return fmt.Errorf("error in starting server: %v", err)
	}

	log.Info("server listen :", slog.String("addres", l.Addr().String()))

	if err := g.GRPCServer.Serve(l); err != nil {
		log.Error("server startup error ", err)
		return fmt.Errorf("error in starting grpc_app server: %v", err)
	}

	return nil
}

func (g *GRPCApp) Stop() {
	g.GRPCServer.GracefulStop()
}
