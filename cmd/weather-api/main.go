package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"weather-api/internal/env"
	"weather-api/internal/handler"
	"weather-api/internal/repository"
	"weather-api/internal/routes"
	"weather-api/internal/service"
)

func main() {
	// get the database URL from environment variables
	migrationURL := env.GetMigrationURL()
	dbURL := env.GetDatabaseURL()

	log.Printf("Trying to run migrations from: %s", migrationURL)

	// run the migrations
	m, err := migrate.New(migrationURL, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v. Migration URL: %s, DB URL: %s", err, migrationURL, dbURL)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrations: No new migrations to apply.")
	} else {
		log.Println("Migrations applied successfully!")
	}

	// connect to the database
	log.Printf("Trying to connect to database: " + dbURL)

	conn, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to Database!")

	// set up the HTTP server
	mux := http.NewServeMux()

	// Weather
	weatherService := service.NewWeatherService()
	weatherHandler := handler.NewWeatherHandler(weatherService)

	// Subscription
	subscriptionRepository := repository.NewSubscriptionRepository(conn)
	subscriptionService := service.NewSubscriptionService(subscriptionRepository)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	// Register routes
	routes.Register(mux, weatherHandler, subscriptionHandler)

	apiPort := env.GetAPIPort()
	server := &http.Server{
		Addr:    ":" + apiPort,
		Handler: mux,
	}

	// Set up graceful shutdown
	done := make(chan bool, 1)

	go gracefulShutdown(done, server, conn)

	// start listening
	log.Printf("Starting server on port %s", apiPort)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}

	// wait for shutdown signal
	<-done
	log.Println("Server shutdown complete")
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
