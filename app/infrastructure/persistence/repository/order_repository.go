package repository

import (
	"context"
	"errors"

	"github.com/mehmetymw/debezium-postgres-es/domain/entity"
	"github.com/mehmetymw/debezium-postgres-es/domain/repository"
	"github.com/mehmetymw/debezium-postgres-es/infrastructure/persistence/models"
	"gorm.io/gorm"
)

// GormOrderRepository implements the OrderRepository interface using GORM
type GormOrderRepository struct {
	db *gorm.DB
}

// NewGormOrderRepository creates a new GormOrderRepository
func NewGormOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &GormOrderRepository{
		db: db,
	}
}

// FindAll retrieves all orders
func (r *GormOrderRepository) FindAll(ctx context.Context) ([]entity.Order, error) {
	var orderModels []models.Order
	if err := r.db.WithContext(ctx).Find(&orderModels).Error; err != nil {
		return nil, err
	}

	orders := make([]entity.Order, len(orderModels))
	for i, model := range orderModels {
		orders[i] = *model.ToEntity()
	}

	return orders, nil
}

// FindByID retrieves an order by its ID
func (r *GormOrderRepository) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	var orderModel models.Order
	if err := r.db.WithContext(ctx).First(&orderModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found
		}
		return nil, err
	}

	return orderModel.ToEntity(), nil
}

// FindByStatus retrieves orders by status
func (r *GormOrderRepository) FindByStatus(ctx context.Context, status string) ([]entity.Order, error) {
	var orderModels []models.Order
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&orderModels).Error; err != nil {
		return nil, err
	}

	orders := make([]entity.Order, len(orderModels))
	for i, model := range orderModels {
		orders[i] = *model.ToEntity()
	}

	return orders, nil
}

// Create creates a new order
func (r *GormOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	orderModel := models.Order{}
	orderModel.FromEntity(order)

	return r.db.WithContext(ctx).Create(&orderModel).Error
}

// Update updates an existing order
func (r *GormOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	orderModel := models.Order{}
	orderModel.FromEntity(order)

	return r.db.WithContext(ctx).Save(&orderModel).Error
}

// Delete deletes an order by its ID
func (r *GormOrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Order{}, "id = ?", id).Error
}
