package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/services/api-gateway/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/services/api-gateway/internal/middleware"
	"github.com/ArteShow/Minecraft-Server-Creator/services/api-gateway/internal/proxy"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeou  = 60 * time.Second
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Port != "" && cfg.Port[0] != ':' {
    	cfg.Port = ":" + cfg.Port
	}

	authRegisterProxy := proxy.NewProxy("http://auth-service:8001", "/auth-service/register")
	authLoginProxy := proxy.NewProxy("http://auth-service:8001", "/auth-service/login")

	handler := http.NewServeMux()
	handler.Handle(
		"/api/"+cfg.APIVersion+"/api-gateway/health",
		middleware.LoggingMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, err := w.Write([]byte("ok"))
				if err != nil {
					http.Error(w, "failed to write status ok", http.StatusInternalServerError)
				}
			}),
		),
	)
	handler.Handle("/register",middleware.LoggingMiddleware(authRegisterProxy))
	handler.Handle("/login", middleware.LoggingMiddleware(authLoginProxy))

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: handler,
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
		log.Println("gateway running on "+cfg.Port)
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
