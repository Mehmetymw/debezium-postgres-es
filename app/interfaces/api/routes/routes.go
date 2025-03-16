package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetymw/debezium-postgres-es/interfaces/api/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, orderHandler *handlers.OrderHandler) {
	// API group
	api := app.Group("/api")

	// Orders routes
	orders := api.Group("/orders")
	orders.Get("/", orderHandler.GetAllOrders)
	orders.Get("/:id", orderHandler.GetOrder)
	orders.Post("/", orderHandler.CreateOrder)
	orders.Put("/:id", orderHandler.UpdateOrder)
	orders.Delete("/:id", orderHandler.DeleteOrder)
	orders.Get("/status/:status", orderHandler.GetOrdersByStatus)

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
}
