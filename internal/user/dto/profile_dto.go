package dto

// UserDTO contains basic user information
type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
}

// UpdatePasswordRequest represents a password update request
type UpdatePasswordRequest struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

type UpdateAvatarRequest struct {
	UserID      string `json:"user_id"`
	FileName    string `json:"file_name"`
	FileContent []byte `json:"file_content"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
}
