package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/smtp"
	"weather-api/internal/config"
	"weather-api/internal/repository"
)

var AlreadySubscribedError = errors.New("already subscribed")
var InvalidTokenError = errors.New("invalid Token")
var TokenNotFoundError = errors.New("token not found")

type SubscriptionService struct {
	repo *repository.Repository
}

func NewSubscriptionService(repo *repository.Repository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
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
		err := ss.repo.UpdateConfirmationToken(token)
		if err != nil {
			log.Println(err)
			return err
		}
		err = sendEmailToMailHog(fromEmail, email, emailSubject, emailBody, token)
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

	err = sendEmailToMailHog(fromEmail, email, emailSubject, emailBody, token)
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

func sendEmailToMailHog(fromEmail, toEmail, subject, body, token string) error {
	smtpHost := config.GetSMTPHost()
	smtpPort := config.GetSMTPPort()
	smtpAddr := smtpHost + ":" + smtpPort

	fullBody := body + "\nYour token: " + token + "\n"
	msgContent := fmt.Sprintf("From: %s \r\nTo: %s \r\nSubject: %s  \r\n\r\n %s", fromEmail, toEmail, subject, fullBody)

	err := smtp.SendMail(smtpAddr, nil, fromEmail, []string{toEmail}, []byte(msgContent))
	if err != nil {
		log.Printf("Failed to send email to %s: %v", toEmail, err)
		return err
	}

	return nil
}
