package api

import (
	"net/http"
	"hot-coffee/internal/service"
)

func NewRouter(
	orderSvc service.OrderService,
	menuSvc service.MenuService,
	inventorySvc service.InventoryService,
	reportsSvc service.ReportsService,
) http.Handler {
	mux := http.NewServeMux()

	// Order endpoints
	mux.HandleFunc("POST /orders", makeHandlerFunc(orderSvc.CreateOrder))
	mux.HandleFunc("GET /orders", makeHandlerFunc(orderSvc.GetOrders))
	mux.HandleFunc("GET /orders/{id}", makeHandlerFunc(orderSvc.GetOrder))
	mux.HandleFunc("PUT /orders/{id}", makeHandlerFunc(orderSvc.UpdateOrder))
	mux.HandleFunc("DELETE /orders/{id}", makeHandlerFunc(orderSvc.DeleteOrder))
	mux.HandleFunc("POST /orders/{id}/close", makeHandlerFunc(orderSvc.CloseOrder))

	// Menu endpoints
	mux.HandleFunc("GET /menu", makeHandlerFunc(menuSvc.GetMenuItems))
	mux.HandleFunc("GET /menu/{id}", makeHandlerFunc(menuSvc.GetMenuItem))

	// Inventory endpoints
	mux.HandleFunc("GET /inventory", makeHandlerFunc(inventorySvc.GetInventoryItems))
	mux.HandleFunc("GET /inventory/{id}", makeHandlerFunc(inventorySvc.GetInventoryItem))
	mux.HandleFunc("PUT /inventory/{id}", makeHandlerFunc(inventorySvc.UpdateInventoryItem))

	// Reports endpoints
	mux.HandleFunc("GET /reports/total-sales", makeHandlerFunc(reportsSvc.GetTotalSales))
	mux.HandleFunc("GET /reports/popular-items", makeHandlerFunc(func() (interface{}, error) {
		return reportsSvc.GetPopularItems(3) // Top 3 popular items
	}))

	return mux
}

type apiHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHandlerFunc[T any](handler func() (T, error)) apiHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		result, err := handler()
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(result)
	}
}

func (f apiHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		// Handle errors appropriately
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}