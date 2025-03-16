package models

import (
	"time"

	"github.com/mehmetymw/debezium-postgres-es/domain/entity"
	"gorm.io/gorm"
)

// Order represents the database model for an order
type Order struct {
	ID         string `gorm:"primaryKey"`
	OrderID    string `gorm:"column:order_id"`
	CustomerID string `gorm:"column:customer_id"`
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for the Order model
func (Order) TableName() string {
	return "orders"
}

// ToEntity converts the model to a domain entity
func (o *Order) ToEntity() *entity.Order {
	return &entity.Order{
		ID:         o.ID,
		OrderID:    o.OrderID,
		CustomerID: o.CustomerID,
		Status:     o.Status,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
		DeletedAt:  o.DeletedAt.Time,
	}
}

// FromEntity converts a domain entity to a model
func (o *Order) FromEntity(order *entity.Order) {
	o.ID = order.ID
	o.OrderID = order.OrderID
	o.CustomerID = order.CustomerID
	o.Status = order.Status
	o.CreatedAt = order.CreatedAt
	o.UpdatedAt = order.UpdatedAt
}
