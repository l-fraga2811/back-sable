package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

// SupabaseAuthMiddleware checks for a valid Supabase JWT token
func SupabaseAuthMiddleware(validator *supabase.TokenValidator) fiber.Handler {
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

		tokenString := parts[1]

		claims, err := validator.Validate(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Set user context
		c.Locals("userID", claims.Subject)
		c.Locals("email", claims.Email)
		c.Locals("token", tokenString)

		return c.Next()
	}
}
