package repository

import (
	"database/sql"
	"errors"
	"log"
)

var AlreadySubscribedError = errors.New("already subscribed")

type SubscriptionRepository struct {
	conn *sql.DB
}

func NewSubscriptionRepository(conn *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		conn: conn,
	}
}

func (repo *SubscriptionRepository) InsertSubscription(email, city, frequency string) error {
	query := `INSERT INTO subs (email, city, frequency, confirm) VALUES ($1, $2, $3, $4)`

	err := repo.GetSubscription(email, city, frequency)
	if err != nil {
		return AlreadySubscribedError
	}

	_, err = repo.conn.Exec(query, email, city, frequency, false)
	if err != nil {
		log.Println("dsadasda" + err.Error())
	}
	return err
}

func (repo *SubscriptionRepository) GetSubscription(email, city, frequency string) error {
	query := `SELECT id FROM subs WHERE email = $1`

	var id int

	err := repo.conn.QueryRow(query, email).Scan(&id)
	if err != nil {
		return err
	}

	log.Println(id)

	return nil
}
