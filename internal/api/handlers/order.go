// internal/api/handlers/order.go
package handlers

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
	"strings"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	switch r.Method {
	case http.MethodPost:
		if len(pathParts) == 1 && pathParts[0] == "orders" {
			h.CreateOrder(w, r)
		} else if len(pathParts) == 3 && pathParts[0] == "orders" && pathParts[2] == "close" {
			h.CloseOrder(w, r, pathParts[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodGet:
		if len(pathParts) == 1 && pathParts[0] == "orders" {
			h.GetOrders(w, r)
		} else if len(pathParts) == 2 && pathParts[0] == "orders" {
			h.GetOrder(w, r, pathParts[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPut:
		if len(pathParts) == 2 && pathParts[0] == "orders" {
			h.UpdateOrder(w, r, pathParts[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodDelete:
		if len(pathParts) == 2 && pathParts[0] == "orders" {
			h.DeleteOrder(w, r, pathParts[1])
		} else {
			http.NotFound(w, r)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.CreateOrder(order)
	if err != nil {
		slog.Error("Failed to create order", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetOrders()
	if err != nil {
		slog.Error("Failed to get orders", "error", err)
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request, id string) {
	order, err := h.service.GetOrder(id)
	if err != nil {
		slog.Error("Failed to get order", "orderID", id, "error", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request, id string) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedOrder, err := h.service.UpdateOrder(id, order)
	if err != nil {
		slog.Error("Failed to update order", "orderID", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteOrder(id); err != nil {
		slog.Error("Failed to delete order", "orderID", id, "error", err)
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request, id string) {
	order, err := h.service.CloseOrder(id)
	if err != nil {
		slog.Error("Failed to close order", "orderID", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
