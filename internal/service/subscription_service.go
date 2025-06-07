package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
	"weather-api/internal/model"
)

var AlreadySubscribedError = errors.New("already subscribed")
var InvalidTokenError = errors.New("invalid Token")
var TokenNotFoundError = errors.New("token not found")

type EmailSender interface {
	Send(from, to, subject, body, token string) error
}

type Subscription interface {
	InsertSubscription(email, city, frequency, token string) error
	GetSubscription(email string) (*model.Subscription, error)
	UpdateTokens(token, unsubscribeToken string) error
	UpdateConfirmationToken(token string) error
	DeleteSubscription(token string) error
}

type SubscriptionService struct {
	repo Subscription
	mail EmailSender
}

func NewSubscriptionService(repo Subscription, mail EmailSender) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
		mail: mail,
	}
}

func (ss *SubscriptionService) Subscribe(email, city, frequency string) error {
	token := uuid.NewString()

	fromEmail := "me@here.com"
	emailSubject := "Subscription Confirmation"
	emailBody := "Please confirm your subscription."

	sub, err := ss.repo.GetSubscription(email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return err
	}

	// if sub exists
	if sub != nil {
		// if the subscription is confirmed
		if sub.Confirmed {
			return AlreadySubscribedError
		}
		err = ss.repo.UpdateConfirmationToken(token)
		if err != nil {
			log.Println(err)
			return err
		}
		err = ss.mail.Send(fromEmail, email, emailSubject, emailBody, token)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}

	err = ss.repo.InsertSubscription(email, city, frequency, token)
	if err != nil {
		log.Println(err)
		return err
	}

	err = ss.mail.Send(fromEmail, email, emailSubject, emailBody, token)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ss *SubscriptionService) Confirm(token string) (string, error) {
	err := uuid.Validate(token)
	if err != nil {
		return "", InvalidTokenError
	}

	unsubscribeToken := uuid.NewString()

	err = ss.repo.UpdateTokens(token, unsubscribeToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", TokenNotFoundError
		}
		log.Println(err)
		return "", err
	}

	return unsubscribeToken, nil
}

func (ss *SubscriptionService) Unsubscribe(token string) error {
	err := uuid.Validate(token)
	if err != nil {
		return InvalidTokenError
	}

	err = ss.repo.DeleteSubscription(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TokenNotFoundError
		}
		return err
	}

	return nil
}
