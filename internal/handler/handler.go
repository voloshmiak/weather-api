package handler

import (
	"net/http"
	"weather-api/internal/service"
)

type Weather interface {
	GetWeather(rw http.ResponseWriter, r *http.Request)
}

type Subscription interface {
	PostSubscription(rw http.ResponseWriter, r *http.Request)
	GetConfirm(rw http.ResponseWriter, r *http.Request)
	GetUnsubscribe(rw http.ResponseWriter, r *http.Request)
}

type Handler struct {
	Weather
	Subscription
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Weather:      NewWeatherHandler(services),
		Subscription: NewSubscriptionHandler(services),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /weather", h.GetWeather)
	mux.HandleFunc("POST /subscribe", h.PostSubscription)
	mux.HandleFunc("GET /confirm/{token}", h.GetConfirm)
	mux.HandleFunc("GET /unsubscribe/{token}", h.GetUnsubscribe)

	mux.Handle("/api/", http.StripPrefix("/api", mux))
}
