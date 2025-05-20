package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"weather-api/internal/config"
	"weather-api/internal/models"
	"weather-api/internal/repository"
)

var CityNotFound = errors.New("city not found")

type WeatherService struct {
	repo *repository.Repository
}

func NewWeatherService(repo *repository.Repository) *WeatherService {
	return &WeatherService{
		repo: repo,
	}
}

func (ws *WeatherService) GetWeatherByCity(city string) (*models.Weather, error) {
	weatherAPIKey := config.GetWeatherAPIKey()
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherAPIKey, city)

	response, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch weather data:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, CityNotFound
	}

	decoder := json.NewDecoder(response.Body)
	var weatherData map[string]interface{}

	if err := decoder.Decode(&weatherData); err != nil {
		return nil, err
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

	return weather, nil
}
