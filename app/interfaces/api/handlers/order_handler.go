package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetymw/debezium-postgres-es/application/service"
	"github.com/mehmetymw/debezium-postgres-es/domain/entity"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// GetAllOrders handles GET /api/orders
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.orderService.GetAllOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching orders",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Orders fetched successfully",
		"data":    orders,
		"count":   len(orders),
	})
}

// GetOrder handles GET /api/orders/:id
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.orderService.GetOrderByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order fetched successfully",
		"data":    order,
	})
}

// CreateOrder handles POST /api/orders
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	order := new(entity.Order)

	// Parse request body
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request",
			"error":   err.Error(),
		})
	}

	// Create order
	if err := h.orderService.CreateOrder(c.Context(), order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating order",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
		"data":    order,
	})
}

// UpdateOrder handles PUT /api/orders/:id
func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	updateData := new(entity.Order)

	// Parse request body
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request",
			"error":   err.Error(),
		})
	}

	// Set ID from path parameter
	updateData.ID = id

	// Update order
	if err := h.orderService.UpdateOrder(c.Context(), updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating order",
			"error":   err.Error(),
		})
	}

	// Get updated order
	updatedOrder, err := h.orderService.GetOrderByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching updated order",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order updated successfully",
		"data":    updatedOrder,
	})
}

// DeleteOrder handles DELETE /api/orders/:id
func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	// Delete order
	if err := h.orderService.DeleteOrder(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting order",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order deleted successfully",
	})
}

// GetOrdersByStatus handles GET /api/orders/status/:status
func (h *OrderHandler) GetOrdersByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	orders, err := h.orderService.GetOrdersByStatus(c.Context(), status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching orders",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Orders fetched successfully",
		"data":    orders,
		"count":   len(orders),
	})
}
