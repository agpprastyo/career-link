package errors

import "github.com/gofiber/fiber/v2"

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends an error response with the specified status code and message
func RespondWithError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(ErrorResponse{
		Error: message,
	})
}
