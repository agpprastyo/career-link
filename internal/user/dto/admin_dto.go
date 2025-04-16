package dto

import "github.com/agpprastyo/career-link/pkg/validator"

type CreateAdminRequest struct {
	Username  string              `json:"username"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	AdminRole string              `json:"admin_role"`
	Validator validator.Validator `json:"-"`
}

type CreateAdminResponse struct {
	Data    DataCreateAdminResponse `json:"data"`
	Message string                  `json:"message"`
}

type DataCreateAdminResponse struct {
	User  UserDTO  `json:"user"`
	Admin AdminDTO `json:"admin"`
}

type AdminDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	AdminRole string `json:"admin_role"`
}

type DeleteAdminResponse struct {
	Message string `json:"message" default:"Admin deleted successfully"`
}
