package handler

import (
	"errors"
	"log"
	"net/http"
	"weather-api/internal/service"
)

type SubscriptionHandler struct {
	services *service.Service
}

func NewSubscriptionHandler(service *service.Service) *SubscriptionHandler {
	return &SubscriptionHandler{
		services: service,
	}
}

func (sh *SubscriptionHandler) PostSubscription(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	city := r.FormValue("city")
	frequency := r.FormValue("frequency")
	log.Println(email, city, frequency)

	if email == "" || city == "" || frequency == "" {
		http.Error(rw, "Invalid input", http.StatusBadRequest)
		return
	}

	if frequency != "daily" && frequency != "hourly" {
		http.Error(rw, "Invalid input", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	err := sh.services.Subscribe(email, city, frequency)
	if err != nil {
		if errors.Is(err, service.AlreadySubscribedError) {
			log.Println(err)
			http.Error(rw, "Already subscribed", 409)
			return
		}
		log.Printf("Failed to subscribe: %v\n", err)
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write([]byte("Subscription successful. Confirmation email sent."))
	if err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (sh *SubscriptionHandler) GetConfirm(rw http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(rw, "Invalid token", http.StatusBadRequest)
		return
	}

	unsubscribeToken, err := sh.services.Confirm(token)
	if err != nil {
		if errors.Is(err, service.InvalidTokenError) {
			http.Error(rw, "Invalid token", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.TokenNotFoundError) {
			http.Error(rw, "Token not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to confirm subscription:", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte("Subscription confirmed successfully. Here is your unsubscribe token: " + unsubscribeToken))
	if err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (sh *SubscriptionHandler) GetUnsubscribe(rw http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(rw, "Invalid token", http.StatusBadRequest)
		return
	}

	err := sh.services.Unsubscribe(token)
	if err != nil {
		if errors.Is(err, service.InvalidTokenError) {
			http.Error(rw, "Invalid token", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.TokenNotFoundError) {
			http.Error(rw, "Token not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to unsubscribe:", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte("Unsubscribed successfully"))
	if err != nil {
		log.Println("Failed to write response:", err)
	}
}
