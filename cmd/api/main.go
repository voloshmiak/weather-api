package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"weather-api/internal/config"
	"weather-api/internal/database/postgres"
	"weather-api/internal/handler"
	"weather-api/internal/mail"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// config initialization
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// connect to the database
	conn, err := postgres.Connect(cfg.DB.User, cfg.DB.Password, cfg.DB.Host,
		cfg.DB.Port, cfg.DB.Name, cfg.DB.MigrationURL())
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("Successfully connected to Database!")

	// set up the HTTP server
	mux := http.NewServeMux()

	// Mail hog email service
	ml := mail.NewMailHog(cfg.Mail.Host, cfg.Mail.Port)

	// Weather
	ws := service.NewWeatherService(cfg.WeatherAPIKey)
	wh := handler.NewWeatherHandler(ws)
	wh.RegisterRoutes(mux)

	// Subscription
	sr := repository.NewSubscriptionRepository(conn)
	ss := service.NewSubscriptionService(sr, ml)
	sh := handler.NewSubscriptionHandler(ss)
	sh.RegisterRoutes(mux)

	// Set up prefix for API routes
	mux.Handle("/api/", http.StripPrefix("/api", mux))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	// Set up graceful shutdown
	done := make(chan bool, 1)

	go gracefulShutdown(done, server, conn)

	// start listening
	log.Printf("Starting server on port %s", cfg.Server.Port)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// wait for shutdown signal
	<-done
	log.Println("Server shutdown complete")

	return nil
}

func gracefulShutdown(done chan bool, server *http.Server, conn *sql.DB) {
	// wait for interruption
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	err := conn.Close()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to close database: %s", err))
	}

	done <- true
}
