// internal/api/handlers/inventory.go
package handlers

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: svc}
}

func (h *InventoryHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdItem, err := h.service.CreateInventoryItem(item)
	if err != nil {
		slog.Error("Failed to create inventory item", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdItem)
}

func (h *InventoryHandler) GetInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetInventoryItems()
	if err != nil {
		slog.Error("Failed to get inventory items", "error", err)
		http.Error(w, "Failed to get inventory items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	item, err := h.service.GetInventoryItem(id)
	if err != nil {
		slog.Error("Failed to get inventory item", "ingredientID", id, "error", err)
		http.Error(w, "Inventory item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.service.UpdateInventoryItem(id, item)
	if err != nil {
		slog.Error("Failed to update inventory item", "ingredientID", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteInventoryItem(id); err != nil {
		slog.Error("Failed to delete inventory item", "ingredientID", id, "error", err)
		http.Error(w, "Failed to delete inventory item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
