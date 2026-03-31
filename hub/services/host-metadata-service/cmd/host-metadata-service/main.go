package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/handlers"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/server"
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

	grpcLis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterHostMetadataServiceServer(grpcServer, server.NewServer())

	go func() {
		log.Println("gRPC server running on :" + cfg.GRPCPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	if cfg.Port != "" && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/host-metadata-service/health", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/host-metadata-service/create", handlers.CreateHostServerMetaDataEntry)
	mux.HandleFunc("/host-metadata-service/delete", handlers.DeleteHostServerMetadata)
	mux.HandleFunc("/host-metadata-service/get", handlers.GetHostServerMetadata)
	mux.HandleFunc("/host-metadata-service/add", handlers.AddServerToHostServer)

	httpServer := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	go func() {
		log.Println("HTTP server running on " + cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown failed: %v", err)
	}

	log.Println("shutdown complete")
}