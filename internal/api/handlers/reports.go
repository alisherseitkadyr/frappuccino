package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"log/slog"
	"net/http"
)

type ReportsHandler struct {
	service service.ReportsService
}

func NewReportsHandler(svc service.ReportsService) *ReportsHandler {
	return &ReportsHandler{service: svc}
}

func (h *ReportsHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	totalSales, err := h.service.GetTotalSales()
	if err != nil {
		slog.Error("Failed to get total sales", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		TotalSales float64 `json:"total_sales"`
	}{
		TotalSales: totalSales,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *ReportsHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.service.GetPopularItems(3) // Default to top 3
	if err != nil {
		slog.Error("Failed to get popular items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(popularItems); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
