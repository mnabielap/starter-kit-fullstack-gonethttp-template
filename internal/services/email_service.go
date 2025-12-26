package services

import (
	"fmt"
	"log"
	"net/smtp"

	"starter-kit-fullstack-gonethttp-template/config"
)

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	// In development, we just log the email to console (Mock)
	if s.cfg.App.Env == "development" && s.cfg.SMTP.Host == "" {
		log.Println("--- MOCK EMAIL ---")
		log.Printf("To: %s\n", to)
		log.Printf("Subject: %s\n", subject)
		log.Printf("Body: %s\n", body)
		log.Println("------------------")
		return nil
	}

	// SMTP Implementation
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTP.Host, s.cfg.SMTP.Port)
	auth := smtp.PlainAuth("", s.cfg.SMTP.Username, s.cfg.SMTP.Password, s.cfg.SMTP.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	return smtp.SendMail(addr, auth, s.cfg.SMTP.From, []string{to}, msg)
}

func (s *emailService) SendResetPasswordEmail(to, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.App.URL, token)
	text := fmt.Sprintf("Dear user,\n\nTo reset your password, click on this link: %s\n\nIf you did not request any password resets, then ignore this email.", resetURL)
	return s.SendEmail(to, "Reset Password", text)
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.cfg.App.URL, token)
	text := fmt.Sprintf("Dear user,\n\nTo verify your email, click on this link: %s\n\nIf you did not create an account, then ignore this email.", verifyURL)
	return s.SendEmail(to, "Email Verification", text)
}