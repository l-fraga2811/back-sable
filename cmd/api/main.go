package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/l-fraga2811/back-sable/internal/config"
	"github.com/l-fraga2811/back-sable/internal/routes"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sable Backend",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // TODO: Restrict this in production
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Setup Routes
	routes.SetupRoutes(app)

	// Start Server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
