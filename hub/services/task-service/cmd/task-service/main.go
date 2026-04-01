package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/handlers"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 5 * time.Minute
	idleTimeou   = 5 * time.Minute
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Port != "" && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/task-service/health", func(w http.ResponseWriter, _ *http.Request) {
		_, err = w.Write([]byte("ok"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/task-service/create", handlers.CreateServer)
	mux.HandleFunc("/task-service/start", handlers.StartServer)
	mux.HandleFunc("/task-service/stop", handlers.StopServer)
	mux.HandleFunc("/task-service/delete", handlers.DeleteServer)
	mux.HandleFunc("/task-service/getStats", handlers.GetServerStats)
	mux.HandleFunc("/task-service/getServerStats", handlers.GetServerStats)
	mux.HandleFunc("/task-service/console/ws", handlers.ServerConsoleWS)
	mux.HandleFunc("/task-service/console/command", handlers.SendServerConsoleCommand)

	mux.HandleFunc("/task-service/backup/create", handlers.CreateBackup)
	mux.HandleFunc("/task-service/backup/list", handlers.ListBackups)
	mux.HandleFunc("/task-service/backup/get", handlers.GetBackup)
	mux.HandleFunc("/task-service/backup/delete", handlers.DeleteBackup)

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeou,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Println("server running on :8013")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("graceful shutdown started")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown failed: %v", err)
	}

	log.Println("shutdown complete")
}
