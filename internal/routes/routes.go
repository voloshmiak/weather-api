package routes

import (
	"net/http"
	"weather-api/internal/handler"
)

func Register(mux *http.ServeMux, wh *handler.WeatherHandler, sh *handler.SubscriptionHandler) {
	mux.HandleFunc("GET /weather", wh.GetWeather)
	mux.HandleFunc("POST /subscribe", sh.PostSubscription)
	mux.HandleFunc("GET /confirm/{token}", sh.GetConfirm)
	mux.HandleFunc("GET /unsubscribe/{token}", sh.GetUnsubscribe)

	mux.Handle("/api/", http.StripPrefix("/api", mux))
}
