package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success:   false,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
	})
}

func ValidationErrorResponse(c *fiber.Ctx, errors []string) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Success:   false,
		Message:   "Validation failed",
		Data:      map[string]interface{}{"errors": errors},
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
	})
}