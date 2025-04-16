package mapper

import (
	"github.com/agpprastyo/career-link/internal/user/dto"
	"github.com/agpprastyo/career-link/internal/user/entity"
)

// ToUserDTO converts an entity.User to dto.UserDTO
func ToUserDTO(user *entity.User) dto.UserDTO {
	return dto.UserDTO{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		Verified: user.IsActive,
	}
}

// ToAdminDTO converts an entity.Admin to dto.AdminDTO
func ToAdminDTO(admin *entity.Admin) dto.AdminDTO {
	return dto.AdminDTO{
		ID:        admin.ID.String(),
		UserID:    admin.UserID.String(),
		AdminRole: string(admin.Role),
	}
}
