package entity

import (
	"github.com/google/uuid"
	"time"
)

type Admin struct {
	ID        uuid.UUID         `json:"id" db:"id"`
	UserID    uuid.UUID         `json:"user_id" db:"user_id"`
	Role      AdministratorRole `json:"role" db:"role"`
	LastLogin *time.Time        `db:"last_login" json:"last_login,omitempty"`
	IsActive  bool              `db:"is_active" json:"is_active"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time        `db:"deleted_at" json:"deleted_at,omitempty"`
}

type AdministratorRole string

const (
	AdminRoleSuper  AdministratorRole = "super"
	AdminRoleAdmin  AdministratorRole = "admin"
	AdminRoleViewer AdministratorRole = "viewer"
)
