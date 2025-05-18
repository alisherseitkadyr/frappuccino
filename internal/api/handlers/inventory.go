package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: svc}
}

func (h *InventoryHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var itemReq struct {
		Name     string  `json:"name"`
		Quantity float64 `json:"quantity"`
		Unit     string  `json:"unit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&itemReq); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdItem, err := h.service.CreateInventoryItem(
		itemReq.Name,
		itemReq.Quantity,
		itemReq.Unit,
	)
	if err != nil {
		slog.Error("Failed to create inventory item", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdItem); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) GetInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetInventoryItems()
	if err != nil {
		slog.Error("Failed to get inventory items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/inventory/{")
	n = strings.TrimSuffix(n, "}")
	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetInventoryItem(id)
	if err != nil {
		slog.Error("Failed to get inventory item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/inventory/{")
	n = strings.TrimSuffix(n, "}")
	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the ID in the path matches the ID in the body
	if item.IngredientID != id {
		http.Error(w, "ID in path doesn't match ID in body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.service.UpdateInventoryItem(item)
	if err != nil {
		slog.Error("Failed to update inventory item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedItem); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/inventory/{")
	n = strings.TrimSuffix(n, "}")
	id, err := strconv.ParseInt(n, 10, 64)

	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteInventoryItem(id); err != nil {
		slog.Error("Failed to delete inventory item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
