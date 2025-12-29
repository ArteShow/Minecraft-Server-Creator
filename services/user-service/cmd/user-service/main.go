package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	"github.com/ArteShow/Minecraft-Server-Creator/user-service/internal/config"
	userpb "github.com/ArteShow/Minecraft-Server-Creator/user-service/internal/proto"
	internalgrpc "github.com/ArteShow/Minecraft-Server-Creator/user-service/internal/server"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	grpcLis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, internalgrpc.NewServer())

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	go func() {
		log.Println("gRPC server running on :50051")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	go func() {
		log.Println("HTTP server running on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	grpcServer.GracefulStop()

	log.Println("Servers stopped")
}
