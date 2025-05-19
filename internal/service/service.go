package service

import (
	"weather-api/internal/models"
	"weather-api/internal/repository"
)

type Weather interface {
	GetWeatherByCity(city string) (*models.Weather, error)
}

type Subscription interface {
	Subscribe(email, city, frequency string) error
	Confirm(token string) (string, error)
	Unsubscribe(token string) error
}

type Service struct {
	Weather
	Subscription
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Weather:      NewWeatherService(repos),
		Subscription: NewSubscriptionService(repos),
	}
}
