package service

import (
	"weather-api/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (ss *SubscriptionService) Subscribe(email, city, frequency string) error {
	return ss.repo.InsertSubscription(email, city, frequency)
}
