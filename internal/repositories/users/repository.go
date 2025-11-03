package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kuzin57/auth-system/internal/entities"
	"github.com/kuzin57/auth-system/internal/repositories"
)

type Repository struct {
	pgDriver *repositories.PgDriver
}

func NewRepository(pgDriver *repositories.PgDriver) *Repository {
	return &Repository{pgDriver: pgDriver}
}

func (r *Repository) CreateUser(ctx context.Context, user entities.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (email, hashed_password)
		VALUES ($1, $2)
		RETURNING id
	`

	var userID uuid.UUID

	err := r.pgDriver.DB().QueryRowxContext(ctx, query, user.Email, user.HashedPassword).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (r *Repository) GetRegisteredUsers(ctx context.Context) ([]entities.User, error) {
	query := `
		SELECT id, email, created_at, updated_at
		FROM users
	`

	var users []entities.User

	err := r.pgDriver.DB().SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get registered users: %w", err)
	}

	return users, nil
}
