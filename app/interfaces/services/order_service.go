package services

import (
	"context"

	"github.com/mehmetymw/debezium-postgres-es/domain/entity"
)

// OrderService defines the interface for order operations
type OrderService interface {
	// GetAllOrders retrieves all orders
	GetAllOrders(ctx context.Context) ([]entity.Order, error)

	// GetOrderByID retrieves an order by its ID
	GetOrderByID(ctx context.Context, id string) (*entity.Order, error)

	// GetOrdersByStatus retrieves orders by status
	GetOrdersByStatus(ctx context.Context, status string) ([]entity.Order, error)

	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, order *entity.Order) error

	// UpdateOrder updates an existing order
	UpdateOrder(ctx context.Context, order *entity.Order) error

	// DeleteOrder deletes an order by its ID
	DeleteOrder(ctx context.Context, id string) error
}
