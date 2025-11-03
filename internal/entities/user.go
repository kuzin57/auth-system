package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `db:"id"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
