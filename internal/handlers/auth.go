package handlers

import (
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) GetProfile(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	email, _ := c.Locals("email").(string)

	return c.JSON(fiber.Map{
		"id":    userID,
		"email": email,
	})
}
