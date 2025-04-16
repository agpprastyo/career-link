package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Role      Role      `json:"role" db:"role"`
	Avatar    *string   `json:"avatar" db:"avatar"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Role string

const (
	AdminRole     Role = "admin"
	CompanyRole   Role = "company"
	JobSeekerRole Role = "job_seeker"
)
