package notifications

import (
	"fmt"
	"net"
	"net/smtp"

	"github.com/rs/zerolog/log"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SimpleEmail struct {
	To      string
	Subject string
	Body    string
}

type EmailNotifier struct {
	config *SMTPConfig
}

func NewEmailNotifier(config *SMTPConfig) *EmailNotifier {
	return &EmailNotifier{config: config}
}

func (e *EmailNotifier) SendSimpleEmail(email *SimpleEmail) error {
	addr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close smtp connection")
		}
	}()

	client, err := smtp.NewClient(conn, e.config.Host)
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Quit(); err != nil {
			log.Error().Err(err).Msg("failed to close smtp client")
		}
	}()

	if e.config.Username != "" || e.config.Password != "" {
		auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(e.config.From); err != nil {
		return err
	}
	if err := client.Rcpt(email.To); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.config.From, email.To, email.Subject, email.Body,
	)
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	return w.Close()
}

func (e *EmailNotifier) SendLoginNotification(userEmail, userName string) error {
	body := fmt.Sprintf(
		`Hello %s,

You have successfully logged into your account.

If this wasn't you, please contact support immediately.

Best Regards,
The Shop Team`, userName,
	)

	email := &SimpleEmail{
		To:      userEmail,
		Subject: "Login Notification",
		Body:    body,
	}
	return e.SendSimpleEmail(email)
}
