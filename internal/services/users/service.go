package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kuzin57/auth-system/internal/broker"
	"github.com/kuzin57/auth-system/internal/entities"
	"github.com/kuzin57/auth-system/internal/models"
	"github.com/kuzin57/auth-system/internal/repositories/users"
	"github.com/kuzin57/auth-system/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	usersRepo *users.Repository
	broker    *broker.MessageBroker
}

func NewService(usersRepo *users.Repository, broker *broker.MessageBroker) *Service {
	return &Service{usersRepo: usersRepo, broker: broker}
}

func (s *Service) SendRegistrationLink(ctx context.Context, request models.SendRegistrationLinkRequest) error {
	msg := models.SendRegistrationLinkMessage(request)

	err := s.broker.PublishSendRegistrationLinkMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish send registration link message: %w", err)
	}

	return nil
}

func (s *Service) RegisterUser(ctx context.Context, request models.RegisterRequest) (uuid.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate hashed password: %w", err)
	}

	user := entities.User{
		Email:          request.Email,
		HashedPassword: string(hashedPassword),
	}

	userID, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (s *Service) ListUsers(ctx context.Context) (models.ListUsersResponse, error) {
	users, err := s.usersRepo.GetRegisteredUsers(ctx)
	if err != nil {
		return models.ListUsersResponse{}, fmt.Errorf("failed to list users: %w", err)
	}

	response := models.ListUsersResponse{
		Users: utils.Map(users, func(user entities.User) models.User {
			return models.User{
				ID:        user.ID,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
		}),
	}

	return response, nil
}
