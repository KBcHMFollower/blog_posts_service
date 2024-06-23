package grpcapp

import (
	"fmt"
	commentservice "github.com/KBcHMFollower/blog_posts_service/internal/services/comment_service"
	postService "github.com/KBcHMFollower/blog_posts_service/internal/services/post_service"
	"log/slog"
	"net"

	grpcserver "github.com/KBcHMFollower/blog_posts_service/internal/grpc"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log        *slog.Logger
	port       int
	GRPCServer *grpc.Server
}

func New(port int, log *slog.Logger, postService *postService.PostService, commService *commentservice.CommentsService) *GRPCApp {
	cleanGrpcServer := grpc.NewServer()
	grpcserver.RegisterPostServer(cleanGrpcServer, postService)
	grpcserver.RegisterCommentsServer(cleanGrpcServer, commService)

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
