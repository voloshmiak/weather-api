package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"weather-api/internal/models"
)

type WeatherHandler struct{}

func (wh *WeatherHandler) GetWeather(rw http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(rw, "City parameter is required", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=56435458281c4b9fa24164805251305&q=%s&aqi=no", city)

	response, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch weather data:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(rw, "City not found", http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(response.Body)
	var weatherData map[string]interface{}

	if err := decoder.Decode(&weatherData); err != nil {
		log.Println("Failed to decode weather data:", err)
	}

	current, _ := weatherData["current"].(map[string]interface{})
	temperature, _ := current["temp_c"].(float64)
	humidity, _ := current["humidity"].(float64)
	condition, _ := current["condition"].(map[string]interface{})
	description, _ := condition["text"].(string)

	weather := &models.Weather{
		Temperature: temperature,
		Humidity:    humidity,
		Description: description,
	}

	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(weather); err != nil {
		log.Println("Failed to encode weather data:", err)
	}

	rw.WriteHeader(http.StatusOK)
}
