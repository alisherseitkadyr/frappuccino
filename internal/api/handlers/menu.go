package handlers

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"log/slog"
	"net/http"
	"strings"
)

type MenuHandler struct {
	service service.MenuService
}

func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{service: svc}
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
	id := strings.TrimPrefix(r.URL.Path, "/menu/")
	if id == "" {
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