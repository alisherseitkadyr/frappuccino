package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"log"
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
		log.Printf("Failed to get total sales", "error", err)
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
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *ReportsHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.service.GetPopularItems(3) // Default to top 3
	if err != nil {
		log.Printf("Failed to get popular items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(popularItems); err != nil {
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *ReportsHandler) SearchReportHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")
	min := r.URL.Query().Get("minPrice")
	max := r.URL.Query().Get("maxPrice")

	if q == "" {
		http.Error(w, "query param 'q' is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.SearchReport(q, filter, min, max)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)


}

func (h *ReportsHandler) OrderedItemsByPeriodHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	period := query.Get("period")
	month := query.Get("month")
	year := query.Get("year")

	resp, err := h.service.GetOrderedItemsByPeriod(period, month, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}