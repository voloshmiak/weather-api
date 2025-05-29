package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"weather-api/internal/service"
)

type WeatherHandler struct {
	services service.Weather
}

func NewWeatherHandler(service service.Weather) *WeatherHandler {
	return &WeatherHandler{services: service}
}

func (wh *WeatherHandler) GetWeather(rw http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	weather, err := wh.services.GetWeatherByCity(city)
	if err != nil {
		if errors.Is(err, service.CityNotFound) {
			http.Error(rw, "City not found", http.StatusNotFound)
			return
		}
		log.Println(err)
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
