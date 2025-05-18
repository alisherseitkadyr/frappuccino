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

type MenuHandler struct {
	service service.MenuService
}

func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{service: svc}
}

func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdItem, err := h.service.CreateMenuItem(item)
	if err != nil {
		slog.Error("Failed to create menu item", "error", err)
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

func (h *MenuHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetMenuItems()
	if err != nil {
		slog.Error("Failed to get menu items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/menu/{")
	n = strings.TrimSuffix(n, "}")
	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "Menu item ID is required", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetMenuItem(id)
	if err != nil {
		slog.Error("Failed to get menu item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/menu/{")
	n = strings.TrimSuffix(n, "}")

	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Menu item ID is required", http.StatusBadRequest)
		return
	}

	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.service.UpdateMenuItem(id, item)
	if err != nil {
		slog.Error("Failed to update menu item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedItem); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	n := strings.TrimPrefix(r.URL.Path, "/menu/{")
	n = strings.TrimSuffix(n, "}")

	id, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		http.Error(w, "Inventory item ID is required", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Menu item ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteMenuItem(id); err != nil {
		slog.Error("Failed to delete menu item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
