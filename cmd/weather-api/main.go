package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"weather-api/internal/handlers"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var (
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		url      = fmt.Sprintf("host=localhost port=5432 dbname=weather-database user=%s password=%s", user, password)
	)

	conn, err := sql.Open("pgx", url)
	if err != nil {
		panic(err)
	}

	// init mux
	mux := http.NewServeMux()

	// init repo
	subscriptionRepo := repository.NewSubscriptionRepository(conn)

	// init service
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	// init handlers
	weatherHandler := new(handlers.WeatherHandler)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// weather
	mux.HandleFunc("GET /weather", weatherHandler.GetWeather)

	// subscription
	mux.HandleFunc("POST /subscribe", subscriptionHandler.PostSubscription)
	mux.HandleFunc("GET /confirm/{token}", subscriptionHandler.GetConfirm)
	mux.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.GetUnsubscribe)

	mux.Handle("/api/", http.StripPrefix("/api", mux))

	// server configuration
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// start listening
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
