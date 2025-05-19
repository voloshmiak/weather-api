package repository

import (
	"database/sql"
	"weather-api/internal/models"
)

type Weather interface {
}

type Subscription interface {
	InsertSubscription(email, city, frequency, token string) error
	GetSubscription(email string) (*models.Subscription, error)
	UpdateTokens(token, unsubscribeToken string) error
	UpdateConfirmationToken(token string) error
	DeleteSubscription(token string) error
}

type Repository struct {
	Subscription
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{
		Subscription: NewSubscriptionRepository(conn),
	}
}
