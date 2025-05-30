package models

import (
	"time"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusProcessing OrderStatus = "processing"
	StatusClosed     OrderStatus = "closed"
	StatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID           int64
	CustomerName string
	Items        []OrderItem
	TotalPrice   float64
	Status       OrderStatus
	CreatedAt    time.Time
	UpdatedAt	 time.Time
}

type OrderItem struct {
	ProductID int64
	ProductName string
	Quantity  int
}