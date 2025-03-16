package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents an order in the system
type Order struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	OrderID    string         `json:"orderId"`
	CustomerID string         `json:"customerId"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// OrderStatus represents the possible statuses of an order
var OrderStatus = struct {
	New         string
	Processing  string
	Completed   string
	Shipped     string
	Delivered   string
	Cancelled   string
	Returned    string
	Pending     string
	OnHold      string
	Backordered string
}{
	New:         "NEW",
	Processing:  "PROCESSING",
	Completed:   "COMPLETED",
	Shipped:     "SHIPPED",
	Delivered:   "DELIVERED",
	Cancelled:   "CANCELLED",
	Returned:    "RETURNED",
	Pending:     "PENDING",
	OnHold:      "ON_HOLD",
	Backordered: "BACKORDERED",
}
