package entity

import (
	"github.com/google/uuid"
	"time"
)

type VerificationToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	Type      string     `json:"type" db:"type"`
	ExpiredAt time.Time  `json:"expired_at" db:"expired_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
}

type TokenType string

const (
	EmailVerification TokenType = "email_verification"
	PasswordReset     TokenType = "password_reset"
)
