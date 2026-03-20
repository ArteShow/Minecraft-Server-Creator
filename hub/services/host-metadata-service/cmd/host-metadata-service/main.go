package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/server"
	"google.golang.org/grpc"
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

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		log.Println("gRPC server running on :" + cfg.GRPCPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}
