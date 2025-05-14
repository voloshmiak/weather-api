package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather-api/internal/models"
)

type WeatherHandler struct{}

func (wh *WeatherHandler) GetWeather(writer http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=56435458281c4b9fa24164805251305&q=%s&aqi=no", city)

	writer.Header().Set("Content-Type", "application/json")

	response, err := http.Get(url)
	if err != nil {
		http.Error(writer, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(writer, "City not found", http.StatusNotFound)
		return
	}

	data := json.NewDecoder(response.Body)
	var weatherData map[string]interface{}
	if err := data.Decode(&weatherData); err != nil {
		http.Error(writer, "Failed to decode weather data", http.StatusInternalServerError)
		return
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

	writer.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(writer)
	encoder.Encode(weather)
}
