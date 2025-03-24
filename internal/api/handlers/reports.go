// internal/api/handlers/reports.go
package handlers

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

type ReportsHandler struct {
	service service.ReportsService
}

func NewReportsHandler(svc service.ReportsService) *ReportsHandler {
	return &ReportsHandler{service: svc}
}

func (h *ReportsHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	total, err := h.service.GetTotalSales()
	if err != nil {
		slog.Error("Failed to calculate total sales", "error", err)
		http.Error(w, "Failed to calculate total sales", http.StatusInternalServerError)
		return
	}

	response := map[string]float64{"total_sales": total}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ReportsHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 0

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			slog.Error("Invalid limit parameter", "limit", limitStr, "error", err)
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	items, err := h.service.GetPopularItems(limit)
	if err != nil {
		slog.Error("Failed to get popular items", "error", err)
		http.Error(w, "Failed to get popular items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
