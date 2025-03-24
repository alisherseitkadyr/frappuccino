package api

import (
	"net/http"
	"hot-coffee/internal/service"
	"hot-coffee/internal/api/handlers"
)

func NewRouter(
	orderSvc service.OrderService,
	menuSvc service.MenuService,
	inventorySvc service.InventoryService,
	reportsSvc service.ReportsService,
) http.Handler {
	mux := http.NewServeMux()

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderSvc)
	menuHandler := handlers.NewMenuHandler(menuSvc)
	inventoryHandler := handlers.NewInventoryHandler(inventorySvc)
	reportsHandler := handlers.NewReportsHandler(reportsSvc)

	// Order endpoints
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrder)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.UpdateOrder)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrder)
	mux.HandleFunc("POST /orders/{id}/close", orderHandler.CloseOrder)

	// Menu endpoints
	mux.HandleFunc("GET /menu", menuHandler.GetMenuItems)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetMenuItem)

	// Inventory endpoints
	mux.HandleFunc("GET /inventory", inventoryHandler.GetInventoryItems)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.GetInventoryItem)
	mux.HandleFunc("PUT /inventory/{id}", inventoryHandler.UpdateInventoryItem)

	// Reports endpoints
	mux.HandleFunc("GET /reports/total-sales", reportsHandler.GetTotalSales)
	mux.HandleFunc("GET /reports/popular-items", reportsHandler.GetPopularItems)

	return mux
}