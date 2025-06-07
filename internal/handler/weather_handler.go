package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"weather-api/internal/config"
	"weather-api/internal/model"
	"weather-api/internal/service"
)

type Weather interface {
	GetWeatherByCity(city, weatherAPIKey string) (*model.Weather, error)
}

type WeatherHandler struct {
	service Weather
	config  *config.Config
}

func NewWeatherHandler(service Weather, config *config.Config) *WeatherHandler {
	return &WeatherHandler{service: service, config: config}
}

func (wh *WeatherHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /weather", wh.GetWeather)
}

func (wh *WeatherHandler) GetWeather(rw http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	weather, err := wh.service.GetWeatherByCity(city, wh.config.WeatherAPIKey)
	if err != nil {
		if errors.Is(err, service.CityNotFound) {
			http.Error(rw, "City not found", http.StatusNotFound)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(weather); err != nil {
		log.Println("Failed to encode weather data:", err)
		return
	}
}
