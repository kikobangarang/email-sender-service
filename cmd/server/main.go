package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/kikobangarang/email-sender-service/internal/api"
	"github.com/kikobangarang/email-sender-service/internal/email"
	"github.com/kikobangarang/email-sender-service/internal/repository"
)

func main() {
	repo, err := repository.NewSQLiteRepository("emails.db")
	if err != nil {
		log.Fatal(err)
	}

	_ = godotenv.Load()

	sender := email.NewSMTPSender(
		mustEnv("SMTP_HOST"),
		mustEnv("SMTP_PORT"),
		mustEnv("SMTP_USER"),
		mustEnv("SMTP_PASS"),
		mustEnv("SMTP_FROM"),
	)

	emailService := email.NewService(*repo)

	workerCfg := email.WorkerConfig{
		WorkerCount:  3,
		PollInterval: 2 * time.Second,
		MaxRetries:   3,
		BatchSize:    10,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workers := email.NewWorkerPool(*repo, *sender, workerCfg)
	workers.Start(ctx)

	mux := http.NewServeMux()
	api.RegisterHandlers(mux, emailService)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      api.LoggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server listening on " + server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	waitForShutdown(server, cancel)
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env var: %s", key)
	}
	return v
}

func waitForShutdown(server *http.Server, cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	cancel()

	ctx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelTimeout()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
