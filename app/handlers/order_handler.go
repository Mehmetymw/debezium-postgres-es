package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetymw/debezium-postgres-es/config"
	"github.com/mehmetymw/debezium-postgres-es/models"
)

// GetAllOrders returns all orders
func GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order

	// Get all orders from the database
	result := config.DB.Find(&orders)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching orders",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Orders fetched successfully",
		"data":    orders,
		"count":   len(orders),
	})
}

// GetOrder returns a single order by ID
func GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Find the order in the database
	result := config.DB.First(&order, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order fetched successfully",
		"data":    order,
	})
}

// CreateOrder creates a new order
func CreateOrder(c *fiber.Ctx) error {
	order := new(models.Order)

	// Parse request body
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if order.ID == "" || order.OrderID == "" || order.CustomerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID, OrderID and CustomerID are required",
		})
	}

	// Set default status if not provided
	if order.Status == "" {
		order.Status = models.OrderStatus.New
	}

	// Create the order in the database
	result := config.DB.Create(&order)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating order",
			"error":   result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
		"data":    order,
	})
}

// UpdateOrder updates an existing order
func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Find the order in the database
	result := config.DB.First(&order, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
			"error":   result.Error.Error(),
		})
	}

	// Parse request body
	updateData := new(models.Order)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request",
			"error":   err.Error(),
		})
	}

	// Update fields if provided
	if updateData.OrderID != "" {
		order.OrderID = updateData.OrderID
	}
	if updateData.CustomerID != "" {
		order.CustomerID = updateData.CustomerID
	}
	if updateData.Status != "" {
		order.Status = updateData.Status
	}

	// Save the updated order
	config.DB.Save(&order)

	return c.JSON(fiber.Map{
		"message": "Order updated successfully",
		"data":    order,
	})
}

// DeleteOrder deletes an order
func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Find the order in the database
	result := config.DB.First(&order, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
			"error":   result.Error.Error(),
		})
	}

	// Delete the order
	config.DB.Delete(&order)

	return c.JSON(fiber.Map{
		"message": "Order deleted successfully",
	})
}
