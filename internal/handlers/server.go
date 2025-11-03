package handlers

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/kuzin57/auth-system/internal/config"
	"go.uber.org/fx"
)

type Server struct {
	config                      *config.Config
	sendRegistrationLinkHandler *SendRegistrationLinkHandler
	registrationPageHandler     *RegistrationPageHandler
	registerHandler             *RegisterHandler
	listUsersHandler            *ListUsersHandler
	listener                    net.Listener
}

func NewServer(
	lc fx.Lifecycle,
	config *config.Config,
	sendRegistrationLinkHandler *SendRegistrationLinkHandler,
	registrationPageHandler *RegistrationPageHandler,
	registerHandler *RegisterHandler,
	listUsersHandler *ListUsersHandler,
) *Server {
	server := &Server{
		config:                      config,
		sendRegistrationLinkHandler: sendRegistrationLinkHandler,
		registrationPageHandler:     registrationPageHandler,
		registerHandler:             registerHandler,
		listUsersHandler:            listUsersHandler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(ctx)
		},
	})

	return server
}

func (s *Server) Start(ctx context.Context) (err error) {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	router.POST("/send-registration-link", s.sendRegistrationLinkHandler.Handle)
	router.POST("/register", s.registerHandler.Handle)
	router.GET("/users", s.listUsersHandler.Handle)
	router.GET("/registration", s.registrationPageHandler.Handle)

	go func() {
		if err := router.RunListener(s.listener); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("stopping server")

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	return nil
}
