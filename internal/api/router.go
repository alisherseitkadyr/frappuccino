package api

import (
	"frappuccino/internal/api/handlers"
	"frappuccino/internal/service"
	"net/http"
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
	mux.HandleFunc("POST /menu", menuHandler.CreateMenuItem)
	mux.HandleFunc("GET /menu", menuHandler.GetMenuItems)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetMenuItem)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.UpdateMenuItem)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteMenuItem)

	// Inventory endpoints
	mux.HandleFunc("POST /inventory", inventoryHandler.CreateInventoryItem)
	mux.HandleFunc("GET /inventory", inventoryHandler.GetInventoryItems)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.GetInventoryItem)
	mux.HandleFunc("PUT /inventory/{id}", inventoryHandler.UpdateInventoryItem)
	mux.HandleFunc("DELETE /inventory/{id}", inventoryHandler.DeleteInventoryItem)

	// Reports endpoints
	mux.HandleFunc("GET /reports/total-sales", reportsHandler.GetTotalSales)
	mux.HandleFunc("GET /reports/popular-items", reportsHandler.GetPopularItems)


	mux.HandleFunc("GET /orders/numberOfOrderedItems", orderHandler.GetNumberOfOrderedItems)
	mux.HandleFunc("GET /reports/search", reportsHandler.SearchReportHandler)
	mux.HandleFunc("GET /reports/orderedItemsByPeriod", reportsHandler.OrderedItemsByPeriodHandler)
	mux.HandleFunc("GET /inventory/getLeftOvers", inventoryHandler.GetLeftOversHandler)


	return mux
}
