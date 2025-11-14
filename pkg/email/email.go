package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(to, subject, body string) error
}

// emailService is the concrete implementation of EmailService
type emailService struct {
	host        string
	port        int
	user        string
	password    string
	senderEmail string
}

// NewEmailService creates a new instance of EmailService
func NewEmailService(host string, port int, user, password, senderEmail string) EmailService {
	return &emailService{
		host:        host,
		port:        port,
		user:        user,
		password:    password,
		senderEmail: senderEmail,
	}
}

// SendEmail sends an email using SMTP
func (s *emailService) SendEmail(to, subject, body string) error {
	// Setup authentication
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	// Compose email headers and body
	headers := make(map[string]string)
	headers["From"] = s.senderEmail
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// Send email
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.senderEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendEmailToMultiple sends an email to multiple recipients
func (s *emailService) SendEmailToMultiple(recipients []string, subject, body string) error {
	// Setup authentication
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	// Compose email headers and body
	headers := make(map[string]string)
	headers["From"] = s.senderEmail
	headers["To"] = strings.Join(recipients, ", ")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// Send email
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.senderEmail, recipients, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
