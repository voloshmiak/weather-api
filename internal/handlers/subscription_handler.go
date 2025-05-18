package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"weather-api/internal/models"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
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

	err := sh.service.Subscribe(email, city, frequency)
	if errors.Is(err, repository.AlreadySubscribedError) {
		http.Error(rw, "Already subscribed", 409)
	}

	rw.WriteHeader(http.StatusOK)

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
