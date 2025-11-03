package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/kuzin57/auth-system/internal/models"
	"github.com/kuzin57/auth-system/internal/services/token"
	"go.uber.org/fx"
)

type Service struct {
	tokenService *token.Service
	smtpHost     string
	smtpPort     int
	smtpPass     string
	fromEmail    string
	linkBaseURL  string
}

func NewService(lc fx.Lifecycle, tokenService *token.Service) *Service {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	service := &Service{
		tokenService: tokenService,
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     port,
		smtpPass:     os.Getenv("SMTP_PASS"),
		fromEmail:    os.Getenv("SMTP_FROM"),
		linkBaseURL:  os.Getenv("REGISTRATION_LINK_BASE"),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return service.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return service.Stop(ctx)
		},
	})

	return service
}

func (s *Service) Start(ctx context.Context) error {
	if s.smtpHost == "" || s.smtpPort == 0 || s.smtpPass == "" || s.fromEmail == "" || s.linkBaseURL == "" {
		return fmt.Errorf(
			"missing required environment variables: smtpHost: %s, smtpPort: %d, smtpPass: %s, fromEmail: %s, linkBaseURL: %s",
			s.smtpHost,
			s.smtpPort,
			s.smtpPass,
			s.fromEmail,
			s.linkBaseURL,
		)
	}

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	log.Println("stopping email service")

	return nil
}

func (s *Service) HandleSendRegistrationLink(ctx context.Context, msg models.SendRegistrationLinkMessage) error {
	log.Println("handling send registration link message", msg)

	token, err := s.tokenService.GenerateToken(msg.Email, time.Minute)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/registration?token=%s", s.linkBaseURL, token)
	subject := "регистрация"
	body := fmt.Sprintf("Ссылка для регистрации: %s\n", link)

	return s.sendEmail(msg.Email, subject, body)
}

func (s *Service) sendEmail(to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	auth := smtp.PlainAuth("", s.fromEmail, s.smtpPass, s.smtpHost)

	msg := "From: " + s.fromEmail + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: 8bit\r\n\r\n" +
		body

	return smtp.SendMail(addr, auth, s.fromEmail, []string{to}, []byte(msg))
}
