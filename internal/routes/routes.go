package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/handlers"
	"github.com/l-fraga2811/back-sable/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Health Check
	healthHandler := handlers.NewHealthHandler()
	api.Get("/health", healthHandler.Check)

	// Protected Routes Group
	// specific handlers will go here
	protected := api.Group("/v1", middleware.SupabaseAuthMiddleware())
	
	protected.Get("/profile", func(c fiber.Ctx) error {
		// Example protected route
		return c.JSON(fiber.Map{
			"message": "Access granted to protected resource",
			"token": c.Locals("user_token"),
		})
	})
}
