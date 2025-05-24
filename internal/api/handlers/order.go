package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Printf("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.CreateOrder(order)
	if err != nil {
		log.Printf("Failed to create order", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdOrder); err != nil {
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetOrders()
	if err != nil {
		log.Printf("Failed to get orders", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/orders/{")
	n = strings.TrimSuffix(n, "}")

	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(id)
	if err != nil {
		log.Printf("Failed to get order", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/orders/{")
	n = strings.TrimSuffix(n, "}")

	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Printf("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedOrder, err := h.service.UpdateOrder(id, order)
	if err != nil {
		log.Printf("Failed to update order", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedOrder); err != nil {
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/orders/{")
	n = strings.TrimSuffix(n, "}")

	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteOrder(id); err != nil {
		log.Printf("Failed to delete order", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/orders/{")
	n = strings.TrimSuffix(n, "}/close")

	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.CloseOrder(id)
	if err != nil {
		log.Printf("Failed to close order", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
