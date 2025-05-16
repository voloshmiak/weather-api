package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"weather-api/internal/models"
)

type SubscriptionHandler struct {
	conn *sql.DB
}

func NewSubscriptionHandler(conn *sql.DB) *SubscriptionHandler {
	return &SubscriptionHandler{
		conn: conn,
	}
}

func (sh *SubscriptionHandler) PostSubscription(rw http.ResponseWriter, r *http.Request) {
	var subscription models.Subscription
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subscription); err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	email := subscription.Email
	city := subscription.City
	frequency := subscription.Frequency

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("%s, %s, %s", email, city, frequency)))
}

func (sh *SubscriptionHandler) GetConfirm(rw http.ResponseWriter, r *http.Request) {
	// Simulate confirming subscription
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Subscription confirmed"))
}

func (sh *SubscriptionHandler) GetUnsubscribe(rw http.ResponseWriter, r *http.Request) {
	// Simulate unsubscribing
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Unsubscribed successfully"))
}
