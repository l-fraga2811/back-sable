package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

// SupabaseAuthMiddleware checks for a valid Supabase JWT token
// This is a placeholder implementation. You will need to integrate the actual Supabase client
// or a JWT validator to verify the token signature against your Supabase project secret.
func SupabaseAuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		token := parts[1]

		// TODO: Validate 'token' using Supabase Go client or JWT library
		// user, err := supabaseClient.Auth.User(token)
		// if err != nil { ... }

		// For now, we'll just pass the token to the context for valid requests
		// In a real app, you would set the User ID or object here
		c.Locals("user_token", token)

		return c.Next()
	}
}
