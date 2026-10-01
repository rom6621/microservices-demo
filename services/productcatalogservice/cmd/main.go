package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/rom6621/microservices-demo/services/productcatalogservice/genproto"
	"github.com/rom6621/microservices-demo/services/productcatalogservice/internal/handler"
	"github.com/rom6621/microservices-demo/services/productcatalogservice/internal/repository"
	"github.com/rom6621/microservices-demo/services/productcatalogservice/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	port := "3550"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("failed to listen", "port", port, "error", err)
		os.Exit(1)
	}

	// Repository
	pr, err := repository.NewProductRepository("products.json")
	if err != nil {
		slog.Error("failed to construct product repository", "error", err)
		os.Exit(1)
	}
	// Usecase
	pu := usecase.NewProductUsecase(pr)
	// Handler
	ph := handler.NewProductHandler(pu)
	hh := health.NewServer()

	srv := grpc.NewServer()
	pb.RegisterProductCatalogServiceServer(srv, ph)
	healthpb.RegisterHealthServer(srv, hh)

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
