package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l-fraga2811/back-sable/internal/handlers"
	"github.com/l-fraga2811/back-sable/internal/middleware"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

func SetupRoutes(app *fiber.App, validator *supabase.TokenValidator, client *supabase.Client) {
	api := app.Group("/api")

	// Health Check
	healthHandler := handlers.NewHealthHandler()
	api.Get("/health", healthHandler.Check)

	// Auth Routes (Public)
	authHandler := handlers.NewAuthHandler(client)
	api.Post("/auth/signin", authHandler.Login)
	api.Post("/auth/signup", authHandler.Register)

	// Protected Routes Group
	protected := api.Group("/", middleware.SupabaseAuthMiddleware(validator))

	// Auth Handler (Protected)
	protected.Get("/auth/profile", authHandler.GetProfile)

	// Item Routes
	itemHandler := handlers.NewItemHandler(client)
	items := protected.Group("/items")
	items.Get("/", itemHandler.GetAll)
	items.Get("/:id", itemHandler.GetByID)
	items.Post("/", itemHandler.Create)
	items.Put("/:id", itemHandler.Update)
	items.Delete("/:id", itemHandler.Delete)
}
