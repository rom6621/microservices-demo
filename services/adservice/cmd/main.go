package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/rom6621/microservice-demo/services/adservice/genproto"
	"github.com/rom6621/microservice-demo/services/adservice/internal/handler"
	"github.com/rom6621/microservice-demo/services/adservice/internal/repository"
	"github.com/rom6621/microservice-demo/services/adservice/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	port := "9555"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("failed to listen", "port", port, "error", err)
		os.Exit(1)
	}

	// Repository
	ar, err := repository.NewAdRepository()
	if err != nil {
		slog.Error("failed to construct repository", "error", err)
		os.Exit(1)
	}
	// Usecase
	au := usecase.NewAdUsecase(ar)
	// Handler
	ah := handler.NewAdHandler(au)
	hh := health.NewServer()

	srv := grpc.NewServer()
	healthpb.RegisterHealthServer(srv, hh)
	pb.RegisterAdServiceServer(srv, ah)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()

		hh.Shutdown()
		srv.GracefulStop()
	}()

	slog.Info("starting grpc server", "port", port)
	if err := srv.Serve(listener); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
