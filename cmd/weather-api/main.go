package main

import (
	"net/http"
	"weather-api/internal/handlers"
)

func main() {
	// init mux
	mux := http.NewServeMux()

	// init handlers
	weatherHandler := new(handlers.WeatherHandler)
	subscriptionHandler := new(handlers.SubscriptionHandler)

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
