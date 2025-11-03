package models

import (
	"time"

	"github.com/google/uuid"
)

type ListUsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
