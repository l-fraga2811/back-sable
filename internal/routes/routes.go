// internal/routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/handlers"
	"github.com/l-fraga2811/back-sable/internal/middleware"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

func SetupRoutes(app *fiber.App, tokenValidator *supabase.TokenValidator, itemHandler *handlers.ItemHandler, authHandler *handlers.AuthHandler, healthHandler *handlers.HealthHandler) {
	api := app.Group("/api")

	// Auth routes - usando funções globais
	auth := api.Group("/auth")
	auth.Post("/signin", handlers.SignIn)
	auth.Post("/signup", handlers.SignUp)
	auth.Get("/profile", middleware.SupabaseAuthMiddleware(tokenValidator), handlers.GetProfile)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.SupabaseAuthMiddleware(tokenValidator))

	// Item routes
	items := protected.Group("/items")
	items.Get("/", itemHandler.GetAll)
	items.Post("/", itemHandler.Create)
	items.Get("/:id", itemHandler.GetByID)
	items.Put("/:id", itemHandler.Update)
	items.Delete("/:id", itemHandler.Delete)

	// Health check
	app.Get("/health", healthHandler.Check)
}
