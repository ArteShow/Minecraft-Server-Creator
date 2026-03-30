package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/middleware"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/proxy"
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

	createHostServerProxy := proxy.NewProxy("http://network-service:8011", "/network-service/create")

	createHostServerMetadataProxy := proxy.NewProxy("http://host-metadata-service:8012", "/host-metadata-service/create")
	deleteHostServerMetadataProxy := proxy.NewProxy("http://host-metadata-service:8012", "/host-metadata-service/delete")
	getHostServerMetadataProxy := proxy.NewProxy("http://host-metadata-service:8012", "/host-metadata-service/get")
	addServerToHostProxy := proxy.NewProxy("http://host-metadata-service:8012", "/host-metadata-service/add")

	createServerProxy := proxy.NewProxy("http://task-service:8013", "/task-service/create")
	startServerProxy := proxy.NewProxy("http://task-service:8013", "/task-service/start")
	stopServerProxy := proxy.NewProxy("http://task-service:8013", "/task-service/stop")
	deleteServerProxy := proxy.NewProxy("http://task-service:8013", "/task-service/delete")
	getServerStatsProxy := proxy.NewProxy("http://task-service:8013", "/task-service/getStats")

	registerUserProxy := proxy.NewProxy("http://auth-service:8014", "/auth-service/user/register")
	loginUserProxy := proxy.NewProxy("http://auth-service:8014", "/auth-service/user/login")

	getBundlekeyProxy := proxy.NewProxy("http://bundle-service:8015", "/bundle-service/create")
	addBundleProxy := proxy.NewProxy("http://bundle-service:8015", "/bundle-service/add")

	handler := http.NewServeMux()
	handler.Handle(
		"/api/"+cfg.APIVersion+"/gateway/health",
		middleware.LoggingMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, err := w.Write([]byte("ok"))
				if err != nil {
					http.Error(w, "failed to write status ok", http.StatusInternalServerError)
				}
			}),
		),
	)
	handler.Handle("/api/"+cfg.APIVersion+"/network/create", middleware.LoggingMiddleware(middleware.AdminAuthMiddleware(createHostServerProxy)))

	handler.Handle("/api/"+cfg.APIVersion+"/host-metadata/create", middleware.LoggingMiddleware(middleware.AdminAuthMiddleware(createHostServerMetadataProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/host-metadata/delete", middleware.LoggingMiddleware(middleware.AdminAuthMiddleware(deleteHostServerMetadataProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/host-metadata/get", middleware.LoggingMiddleware(middleware.AdminAuthMiddleware(getHostServerMetadataProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/host-metadata/add", middleware.LoggingMiddleware(middleware.AdminAuthMiddleware(addServerToHostProxy)))

	handler.Handle("/api/"+cfg.APIVersion+"/server/create", middleware.LoggingMiddleware(middleware.AuthMiddleware(createServerProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/server/start", middleware.LoggingMiddleware(middleware.AuthMiddleware(startServerProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/server/stop", middleware.LoggingMiddleware(middleware.AuthMiddleware(stopServerProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/server/delete", middleware.LoggingMiddleware(middleware.AuthMiddleware(deleteServerProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/server/getStats", middleware.LoggingMiddleware(middleware.AuthMiddleware(getServerStatsProxy)))

	handler.Handle("/api/"+cfg.APIVersion+"/auth/user/register", middleware.LoggingMiddleware(registerUserProxy))
	handler.Handle("/api/"+cfg.APIVersion+"/auth/user/login", middleware.LoggingMiddleware(loginUserProxy))
	
	handler.Handle("/api/"+cfg.APIVersion+"/bundle/create", middleware.LoggingMiddleware(middleware.AuthMiddleware(getBundlekeyProxy)))
	handler.Handle("/api/"+cfg.APIVersion+"/bundle/add", middleware.LoggingMiddleware(middleware.AuthMiddleware(addBundleProxy)))

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
