package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/handlers"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/server"
	"google.golang.org/grpc"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 60 * time.Second
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	httpPort := cfg.Port
	if httpPort != "" && httpPort[0] != ':' {
		httpPort = ":" + httpPort
	}

	grpcPort := cfg.GRPCPort
	if grpcPort != "" && grpcPort[0] != ':' {
		grpcPort = ":" + grpcPort
	}

	grpcLis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterBundleServiceServer(grpcServer, server.NewServer())

	mux := http.NewServeMux()

	mux.HandleFunc("/bundle-service/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/bundle-service/create", handlers.CreateBundle)

	httpServer := &http.Server{
		Addr:         httpPort,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("HTTP server running on", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	go func() {
		log.Println("gRPC server running on", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("grpc server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down...")

	grpcServer.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown failed: %v", err)
	}

	log.Println("shutdown complete")
}