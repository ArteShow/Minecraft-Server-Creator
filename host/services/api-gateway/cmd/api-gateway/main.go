package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/host/services/api-gateway/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/api-gateway/internal/middleware"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/api-gateway/internal/proxy"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeou   = 60 * time.Second
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Port != "" && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	createServerProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/create")
	startServerProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/start")
	stopServerProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/stop")
	deleteServerProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/delete")
	getServerStatsProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/getServerStats")
	getServerConsoleProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/console")
	sendServerConsoleCommandProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/console/command")

	createBackupProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/backup/create")
	getBackupProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/backup/get")
	deleteBackupProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/backup/delete")
	uploadBackupProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/backup/upload")

	uploadWorldProxy := proxy.NewProxy("http://server-service-v2:8003", "/server-service/world/upload")

	handler := http.NewServeMux()
	handler.Handle(
		"/api-gateway/health",
		middleware.LoggingMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, err := w.Write([]byte("ok"))
				if err != nil {
					http.Error(w, "failed to write status ok", http.StatusInternalServerError)
				}
			}),
		),
	)
	handler.Handle("/server/create", middleware.LoggingMiddleware(middleware.AuthMiddleware()(createServerProxy)))
	handler.Handle("/server/start", middleware.LoggingMiddleware(middleware.AuthMiddleware()(startServerProxy)))
	handler.Handle("/server/stop", middleware.LoggingMiddleware(middleware.AuthMiddleware()(stopServerProxy)))
	handler.Handle("/server/delete", middleware.LoggingMiddleware(middleware.AuthMiddleware()(deleteServerProxy)))
	handler.Handle("/server/getServerStats", middleware.LoggingMiddleware(getServerStatsProxy))
	handler.Handle("/server/console", middleware.LoggingMiddleware(middleware.AuthMiddleware()(getServerConsoleProxy)))
	handler.Handle("/server/console/command", middleware.LoggingMiddleware(middleware.AuthMiddleware()(sendServerConsoleCommandProxy)))

	handler.Handle("/server/backup/create", middleware.LoggingMiddleware(createBackupProxy))
	handler.Handle("/server/backup/get", middleware.LoggingMiddleware(getBackupProxy))
	handler.Handle("/server/backup/delete", middleware.LoggingMiddleware(deleteBackupProxy))
	handler.Handle("/server/backup/upload", middleware.LoggingMiddleware(uploadBackupProxy))

	handler.Handle("/server/world/upload", middleware.LoggingMiddleware(uploadWorldProxy))

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      handler,
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
		log.Println("gateway running on " + cfg.Port)
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
