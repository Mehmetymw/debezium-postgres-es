package service

import (
	"context"
	"errors"
	"time"

	"github.com/mehmetymw/debezium-postgres-es/domain/entity"
	"github.com/mehmetymw/debezium-postgres-es/domain/repository"
)

// OrderService defines the service for order operations
type OrderService struct {
	orderRepo repository.OrderRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	return s.orderRepo.FindAll(ctx)
}

// GetOrderByID retrieves an order by its ID
func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

// GetOrdersByStatus retrieves orders by status
func (s *OrderService) GetOrdersByStatus(ctx context.Context, status string) ([]entity.Order, error) {
	return s.orderRepo.FindByStatus(ctx, status)
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, order *entity.Order) error {
	// Set timestamps
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// Validate order
	if order.ID == "" || order.OrderID == "" || order.CustomerID == "" {
		return errors.New("id, orderId and customerId are required")
	}

	// Set default status if not provided
	if order.Status == "" {
		order.Status = entity.OrderStatus.New
	}

	return s.orderRepo.Create(ctx, order)
}

// UpdateOrder updates an existing order
func (s *OrderService) UpdateOrder(ctx context.Context, order *entity.Order) error {
	// Check if order exists
	existingOrder, err := s.orderRepo.FindByID(ctx, order.ID)
	if err != nil {
		return err
	}
	if existingOrder == nil {
		return errors.New("order not found")
	}

	// Update only provided fields
	if order.OrderID != "" {
		existingOrder.OrderID = order.OrderID
	}
	if order.CustomerID != "" {
		existingOrder.CustomerID = order.CustomerID
	}
	if order.Status != "" {
		existingOrder.Status = order.Status
	}

	// Update timestamp
	existingOrder.UpdatedAt = time.Now()

	return s.orderRepo.Update(ctx, existingOrder)
}

// DeleteOrder deletes an order by its ID
func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	// Check if order exists
	existingOrder, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingOrder == nil {
		return errors.New("order not found")
	}

	return s.orderRepo.Delete(ctx, id)
}
