package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"weather-api/internal/model"
)

var CityNotFound = errors.New("city not found")
var NoResponseError = errors.New("no response from weather API")

type WeatherService struct{}

func NewWeatherService() *WeatherService {
	return new(WeatherService)
}

func (ws *WeatherService) GetWeatherByCity(city, weatherAPIKey string) (*model.Weather, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		weatherAPIKey, city)

	response, err := http.Get(url)
	if err != nil {
		return nil, NoResponseError
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusForbidden {
		return nil, errors.New("probably invalid API key")
	}

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

	weather := &model.Weather{
		Temperature: temperature,
		Humidity:    humidity,
		Description: description,
	}

	return weather, nil
}
