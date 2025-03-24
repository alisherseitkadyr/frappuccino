// internal/api/handlers/menu.go
package handlers

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	json.NewEncoder(w).Encode(createdItem)
}

func (h *MenuHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetMenuItems()
	if err != nil {
		slog.Error("Failed to get menu items", "error", err)
		http.Error(w, "Failed to get menu items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	item, err := h.service.GetMenuItem(id)
	if err != nil {
		slog.Error("Failed to get menu item", "productID", id, "error", err)
		http.Error(w, "Menu item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.service.UpdateMenuItem(id, item)
	if err != nil {
		slog.Error("Failed to update menu item", "productID", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteMenuItem(id); err != nil {
		slog.Error("Failed to delete menu item", "productID", id, "error", err)
		http.Error(w, "Failed to delete menu item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
