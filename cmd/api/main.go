// cmd/api/main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/l-fraga2811/back-sable/internal/config"
	"github.com/l-fraga2811/back-sable/internal/handlers"
	"github.com/l-fraga2811/back-sable/internal/repository"
	"github.com/l-fraga2811/back-sable/internal/repository/supabase"
	"github.com/l-fraga2811/back-sable/internal/routes"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Dependencies
	tokenValidator := supabase.NewTokenValidator(cfg)
	supabaseClient := supabase.NewClient(cfg)

	// Initialize Global Auth Handlers
	handlers.InitAuthHandlers(supabaseClient)

	// Initialize GORM Repository
	itemRepo := repository.NewItemRepositoryGORM(cfg.DB)

	// Initialize Handlers
	itemHandler := handlers.NewItemHandler(itemRepo)
	healthHandler := handlers.NewHealthHandler()

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sable Backend",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Setup Routes
	routes.SetupRoutes(app, tokenValidator, itemHandler, nil, healthHandler)

	// Start Server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
