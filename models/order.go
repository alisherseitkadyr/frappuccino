package models

import (
	"time"
)

type Order struct {
	ID           string      `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderRequest struct {
	CustomerName string             `json:"customer_name"`
	Items        []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func NewOrder(customerName string, items []OrderItem) Order {
	return Order{
		ID:           generateOrderID(),
		CustomerName: customerName,
		Items:        items,
		Status:       "open",
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
}

func generateOrderID() string {
	return "order_" + time.Now().Format("20060102150405")
}
