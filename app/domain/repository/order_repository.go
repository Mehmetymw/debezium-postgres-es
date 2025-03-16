package repository

import (
	"context"

	"github.com/mehmetymw/debezium-postgres-es/domain/entity"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	// FindAll retrieves all orders
	FindAll(ctx context.Context) ([]entity.Order, error)

	// FindByID retrieves an order by its ID
	FindByID(ctx context.Context, id string) (*entity.Order, error)

	// FindByStatus retrieves orders by status
	FindByStatus(ctx context.Context, status string) ([]entity.Order, error)

	// Create creates a new order
	Create(ctx context.Context, order *entity.Order) error

	// Update updates an existing order
	Update(ctx context.Context, order *entity.Order) error

	// Delete deletes an order by its ID
	Delete(ctx context.Context, id string) error
}
