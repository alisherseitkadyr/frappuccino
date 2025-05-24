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
	ID           int64     // ID будет заполнен после вставки в БД
	Customer_name string
	Items        []OrderItem
	TotalPrice   float64
	Status       OrderStatus
	CreatedAt    time.Time
	UpdatedAt	 time.Time
}

type OrderItem struct {
	ProductID int64
	Quantity  int
}

// type OrderRequest struct {
// 	CustomerName string           
// 	Items        []OrderItemRequest 
// }

// type OrderItemRequest struct {
// 	ProductID int64 
// 	Quantity  int   
// }

// NewOrder создаёт объект заказа с начальным статусом "pending" и текущим временем.
// func NewOrder(customerName string, items []OrderItem) Order {
// 	return Order{
// 		ID:           0, // ID присвоит база
// 		customer_name: 0,
// 		Items:        items,
// 		Status:       StatusPending,
// 		CreatedAt:    time.Now(),
// 	}
// }
