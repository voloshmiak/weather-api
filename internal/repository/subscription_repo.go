package repository

import (
	"database/sql"
	"weather-api/internal/model"
)

type SubscriptionRepository struct {
	conn *sql.DB
}

func NewSubscriptionRepository(conn *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		conn: conn,
	}
}

func (repo *SubscriptionRepository) InsertSubscription(email, city, frequency, token string) error {
	query := `INSERT INTO subs (email, city, frequency, confirmed, confirmation_token) VALUES ($1, $2, $3, $4, $5)`

	_, err := repo.conn.Exec(query, email, city, frequency, false, token)
	if err != nil {
		return err
	}
	return err
}

func (repo *SubscriptionRepository) GetSubscription(email string) (*model.Subscription, error) {
	sub := new(model.Subscription)

	query := `SELECT email, city, frequency, confirmed FROM subs WHERE email = $1`

	err := repo.conn.QueryRow(query, email).Scan(&sub.Email, &sub.City, &sub.Frequency, &sub.Confirmed)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (repo *SubscriptionRepository) UpdateTokens(token, unsubscribeToken string) error {
	query := `UPDATE subs SET confirmed = true, unsubscribe_token = $1  WHERE confirmation_token = $2`
	_, err := repo.conn.Exec(query, unsubscribeToken, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SubscriptionRepository) UpdateConfirmationToken(token string) error {
	query := `UPDATE subs SET confirmation_token = $1 WHERE confirmation_token = $2`
	_, err := repo.conn.Exec(query, token, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SubscriptionRepository) DeleteSubscription(token string) error {
	query := `DELETE FROM subs WHERE unsubscribe_token = $1`
	_, err := repo.conn.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}
