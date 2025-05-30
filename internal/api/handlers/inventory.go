package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log"
	"math"
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
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Unit     string `json:"unit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&itemReq); err != nil {
		log.Printf("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdItem, err := h.service.CreateInventoryItem(
		itemReq.Name,
		itemReq.Quantity,
		itemReq.Unit,
	)
	if err != nil {
		log.Printf("Failed to create inventory item", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdItem); err != nil {
		log.Printf("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *InventoryHandler) GetInventoryItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetInventoryItems()
	if err != nil {
		log.Printf("Failed to get inventory items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("Failed to encode response", "error", err)
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
		log.Printf("Failed to get inventory item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		log.Printf("Failed to encode response", "error", err)
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
		log.Printf("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if item.IngredientID != id {
		http.Error(w, "ID in path doesn't match ID in body", http.StatusBadRequest)
		return
	}

	updatedItem, err := h.service.UpdateInventoryItem(item)
	if err != nil {
		log.Printf("Failed to update inventory item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedItem); err != nil {
		log.Printf("Failed to encode response", "error", err)
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
		log.Printf("Failed to delete inventory item", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) GetLeftOversHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	items, total, err := h.service.GetLeftOvers(sortBy, page, pageSize)
	if err != nil {
		http.Error(w, "Error fetching leftovers", http.StatusInternalServerError)
		return
	}

	type ResponseItem struct {
		Name     string  `json:"name"`
		Quantity int `json:"quantity"`
	}
	var responseData []ResponseItem
	for _, item := range items {
		responseData = append(responseData, ResponseItem{
			Name:     item.Name,
			Quantity: item.Quantity,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	hasNext := page < totalPages

	resp := map[string]interface{}{
		"currentPage": page,
		"hasNextPage": hasNext,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
		"data":        responseData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
