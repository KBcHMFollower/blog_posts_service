package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	grpcserver "github.com/KBcHMFollower/test_plate_blog_service/internal/grpc"
	postService "github.com/KBcHMFollower/test_plate_blog_service/internal/services"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log        *slog.Logger
	port       int
	GRPCServer *grpc.Server
}

func New(port int, log *slog.Logger, postService *postService.PostService) *GRPCApp {
	cleanGrpcServer := grpc.NewServer()
	grpcserver.Register(cleanGrpcServer, postService)

	return &GRPCApp{
		log:        log,
		port:       port,
		GRPCServer: cleanGrpcServer,
	}
}

func (g *GRPCApp) Run() error {
	const op = "GRPCApp.Run"
	log := g.log.With(
		slog.String("op", op),
	)

	log.Info("grpc server is trying to get up")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		log.Error("server startup error ", err)
		return fmt.Errorf("error in starting server: %v", err)
	}

	log.Info("server listen :", slog.String("addres", l.Addr().String()))

	if err := g.GRPCServer.Serve(l); err != nil {
		log.Error("server startup error ", err)
		return fmt.Errorf("error in starting grpc server: %v", err)
	}

	return nil
}

func (g *GRPCApp) Stop() {
	const op = "GRPCApp.Stop"
	log := g.log.With(
		slog.String("op", op),
	)

	log.Info("Server is trying to stop")

	g.GRPCServer.GracefulStop()
}
