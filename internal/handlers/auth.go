package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

type AuthHandler struct {
	client *supabase.Client
}

func NewAuthHandler(client *supabase.Client) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var creds supabase.SignInCredentials
	if err := c.Bind().Body(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.client.SignIn(c.Context(), creds)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication failed: " + err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var creds supabase.SignUpCredentials
	if err := c.Bind().Body(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Extract additional fields from request body
	var requestBody map[string]interface{}
	if err := c.Bind().Body(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Map frontend fields to Supabase data structure
	creds.Data = make(map[string]interface{})
	if username, ok := requestBody["username"]; ok {
		creds.Data["username"] = username
	}
	if phone, ok := requestBody["phone"]; ok && phone != "" {
		creds.Data["phone"] = phone
	}

	response, err := h.client.SignUp(c.Context(), creds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Registration failed: " + err.Error(),
		})
	}

	return c.JSON(response)
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
