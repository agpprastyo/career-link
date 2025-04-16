package dto

import (
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/pkg/validator"
	"time"
)

// LoginRequest represents login credentials
// @Description Login request payload
type LoginRequest struct {
	Email     string              `json:"email"`
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Validator validator.Validator `json:"-"`
}

// LoginResponse represents successful login response
// @Description Login response with user data and token
type LoginResponse struct {
	User   entity.User `json:"user"`
	Token  string      `json:"token"`
	Expiry time.Time   `json:"expiry"`
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Username         string               `json:"username"`
	Email            string               `json:"email"`
	Password         string               `json:"password"`
	Role             string               `json:"role"`
	JobSeekerProfile *JobSeekerProfileDTO `json:"job_seeker_profile,omitempty"`
	CompanyProfile   *CompanyProfileDTO   `json:"company_profile,omitempty"`
	Validator        validator.Validator  `json:"-"`
}

// RegisterResponse represents registration response data
type RegisterResponse struct {
	User    entity.User `json:"user"`
	Message string      `json:"message"`
}

// JobSeekerProfileDTO contains profile data for job seekers during registration
type JobSeekerProfileDTO struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Bio         string `json:"bio,omitempty"`
}

// CompanyProfileDTO contains profile data for companies during registration
type CompanyProfileDTO struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Industry    string `json:"industry,omitempty"`
	Website     string `json:"website,omitempty"`
	Email       string `json:"email"`
	Phone       string `json:"phone,omitempty"`
}

type VerifyEmailRequest struct {
	Token               string `json:"token"`
	Email               string `json:"email"`
	validator.Validator `json:"-"`
}

type ResendVerificationRequest struct {
	Email     string              `json:"email"`
	Validator validator.Validator `json:"-"`
}

type LogoutResponse struct {
	Message string `json:"message" default:"User logged out successfully"`
}

// ForgotPasswordRequest represents a password reset request
type ForgotPasswordRequest struct {
	Email     string `json:"email"`
	Validator validator.Validator
}

// ResetPasswordRequest represents a new password submission
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
	Validator   validator.Validator
}
