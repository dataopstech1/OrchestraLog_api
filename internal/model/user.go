package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	FirstName    string     `db:"first_name" json:"first_name"`
	LastName     string     `db:"last_name" json:"last_name"`
	Role         string     `db:"role" json:"role"`
	Department   *string    `db:"department" json:"department"`
	AvatarURL    *string    `db:"avatar_url" json:"avatar_url"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type AuditLog struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     *uuid.UUID `db:"user_id" json:"user_id"`
	Action     string     `db:"action" json:"action"`
	EntityType string     `db:"entity_type" json:"entity_type"`
	EntityID   *string    `db:"entity_id" json:"entity_id"`
	Details    []byte     `db:"details" json:"details,omitempty"`
	IPAddress  *string    `db:"ip_address" json:"ip_address"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}
