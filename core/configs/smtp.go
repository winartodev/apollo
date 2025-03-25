package configs

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	dialer *gomail.Dialer
	Sender string
}

func NewSMTPClient(config SMTP) (*SMTPClient, error) {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Sender, config.Password)
	return &SMTPClient{
		dialer: dialer,
		Sender: config.Sender,
	}, nil
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
	HTML    bool
}

func (c *SMTPClient) Send(email *Email) error {
	message := gomail.NewMessage()

	message.SetHeader("From", c.Sender)
	message.SetHeader("To", email.To)
	message.SetHeader("Subject", email.Subject)

	contentType := "text/plain"
	if email.HTML {
		contentType = "text/html"
	}
	message.SetBody(contentType, email.Body)

	if err := c.dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
