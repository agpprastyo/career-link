package entity

import (
	"github.com/google/uuid"
	"time"
)

// CompanyStatus is an enum-like type for company status
type CompanyStatus string

const (
	Active   CompanyStatus = "active"
	Inactive CompanyStatus = "inactive"
)

// CompanySizeRange  is an enum-like type for company size range
type CompanySizeRange string

const (
	Size1To5      CompanySizeRange = "1-5"
	Size5To10     CompanySizeRange = "5-10"
	Size10To25    CompanySizeRange = "10-25"
	Size25To50    CompanySizeRange = "25-50"
	Size50To100   CompanySizeRange = "50-100"
	Size100To500  CompanySizeRange = "100-500"
	Size500To1000 CompanySizeRange = "500-1000"
	Size1000Plus  CompanySizeRange = "1000+"
)

// Company represents a company entity
type Company struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	UserID      uuid.UUID        `json:"user_id" db:"user_id"` // Foreign key to users table
	Name        string           `json:"name" db:"name"`
	Description *string          `json:"description,omitempty" db:"description"` // Use pointer for nullable fields
	Industry    *string          `json:"industry,omitempty" db:"industry"`
	Website     *string          `json:"website,omitempty" db:"website"`
	Email       string           `json:"email" db:"email"`
	Phone       *string          `json:"phone,omitempty" db:"phone"`
	LogoURL     *string          `json:"logo_url,omitempty" db:"logo_url"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	Status      CompanyStatus    `json:"status" db:"status"`
	Size        CompanySizeRange `json:"size" db:"size"`
	IsVerified  bool             `json:"is_verified" db:"is_verified" default:"false"`
	AddressID   uuid.UUID        `json:"address_id" db:"address_id"` // Foreign key to company_addresses
	Address     *CompanyAddress  `json:"address,omitempty"`          // One-to-one relationship
}

// CompanyAddress represents a company's address entity
type CompanyAddress struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CompanyID    uuid.UUID `json:"company_id" db:"company_id"` // Foreign key to company table
	AddressLine1 string    `json:"address_line_1" db:"address_line_1"`
	AddressLine2 *string   `json:"address_line_2,omitempty" db:"address_line_2"` // Use pointer for nullable fields
	City         string    `json:"city" db:"city"`
	State        string    `json:"state" db:"state"`
	ZipCode      string    `json:"zip_code" db:"zip_code"`
	Country      string    `json:"country" db:"country"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
