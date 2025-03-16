package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mehmetymw/debezium-postgres-es/application/service"
	"github.com/mehmetymw/debezium-postgres-es/config"
	"github.com/mehmetymw/debezium-postgres-es/infrastructure/persistence/migrations"
	"github.com/mehmetymw/debezium-postgres-es/infrastructure/persistence/repository"
	"github.com/mehmetymw/debezium-postgres-es/interfaces/api/handlers"
	"github.com/mehmetymw/debezium-postgres-es/interfaces/api/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	if err := config.ConnectDB(&cfg.PostgreSQL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := migrations.RunMigrations(config.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	orderRepo := repository.NewGormOrderRepository(config.DB)

	// Initialize services
	orderService := service.NewOrderService(orderRepo)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Debezium PostgreSQL to Elasticsearch",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app, orderHandler)

	// Start server
	port := cfg.Server.Port
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
