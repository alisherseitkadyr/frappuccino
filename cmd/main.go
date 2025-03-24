// cmd/main.go
package main

import (
	"flag"
	"fmt"
	"hot-coffee/internal/api"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/service"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 8080, "Port number")
	dataDir := flag.String("dir", "data", "Path to the data directory")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		printUsage()
		return
	}

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(*dataDir)
	menuRepo := repository.NewMenuRepository(*dataDir)
	inventoryRepo := repository.NewInventoryRepository(*dataDir)

	// Initialize services
	orderSvc := service.NewOrderService(orderRepo, menuRepo, inventoryRepo)
	menuSvc := service.NewMenuService(menuRepo)
	inventorySvc := service.NewInventoryService(inventoryRepo)
	reportsSvc := service.NewReportsService(orderRepo, menuRepo)

	// Initialize router
	router := api.NewRouter(orderSvc, menuSvc, inventorySvc, reportsSvc)

	slog.Info("Starting server", "port", *port, "dataDir", *dataDir)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.`)
}
