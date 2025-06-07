package mail

import (
	"fmt"
	"net/smtp"
)

type Hog struct {
	smtpHost string
	smtpPort string
}

func NewMailHog(smtpHost, smtpPort string) *Hog {
	return &Hog{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (e *Hog) Send(from, to, subject, body, token string) error {
	smtpAddr := e.smtpHost + ":" + e.smtpPort
	fullBody := body + "\nYour token: " + token + "\n"
	msgContent := fmt.Sprintf("From: %s \r\nTo: %s \r\nSubject: %s  \r\n\r\n %s", from, to, subject, fullBody)

	err := smtp.SendMail(smtpAddr, nil, from, []string{to}, []byte(msgContent))
	if err != nil {
		return err
	}

	return nil
}
