package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"hot-coffee/internal/api"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		slog.Error("Failed to create data directory", "error", err)
		os.Exit(1)
	}

	// Initialize data files with empty arrays if they don't exist
	if err := initDataFiles(*dataDir); err != nil {
		slog.Error("Failed to initialize data files", "error", err)
		os.Exit(1)
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

	slog.Info("Starting server", "port", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func initDataFiles(dataDir string) error {
	// Initialize orders.json
	ordersFile := filepath.Join(dataDir, "orders.json")
	if _, err := os.Stat(ordersFile); os.IsNotExist(err) {
		if err := writeEmptyArrayToFile(ordersFile); err != nil {
			return fmt.Errorf("failed to initialize orders.json: %w", err)
		}
	}

	// Initialize menu_items.json with sample data if empty
	menuFile := filepath.Join(dataDir, "menu_items.json")
	if _, err := os.Stat(menuFile); os.IsNotExist(err) {
		sampleMenu := []models.MenuItem{
			{
				ID:          "espresso",
				Name:        "Espresso",
				Description: "Strong coffee shot",
				Price:       2.50,
				Ingredients: []models.MenuItemIngredient{
					{IngredientID: "coffee_beans", Quantity: 10},
					{IngredientID: "water", Quantity: 30},
				},
			},
			{
				ID:          "latte",
				Name:        "Latte",
				Description: "Espresso with steamed milk",
				Price:       3.50,
				Ingredients: []models.MenuItemIngredient{
					{IngredientID: "coffee_beans", Quantity: 10},
					{IngredientID: "water", Quantity: 30},
					{IngredientID: "milk", Quantity: 200},
				},
			},
		}
		if err := writeToFile(menuFile, sampleMenu); err != nil {
			return fmt.Errorf("failed to initialize menu_items.json: %w", err)
		}
	}

	// Initialize inventory.json with sample data if empty
	inventoryFile := filepath.Join(dataDir, "inventory.json")
	if _, err := os.Stat(inventoryFile); os.IsNotExist(err) {
		sampleInventory := []models.InventoryItem{
			{
				IngredientID: "coffee_beans",
				Name:         "Coffee Beans",
				Quantity:     1000,
				Unit:         "g",
			},
			{
				IngredientID: "water",
				Name:         "Water",
				Quantity:     5000,
				Unit:         "ml",
			},
			{
				IngredientID: "milk",
				Name:         "Milk",
				Quantity:     3000,
				Unit:         "ml",
			},
		}
		if err := writeToFile(inventoryFile, sampleInventory); err != nil {
			return fmt.Errorf("failed to initialize inventory.json: %w", err)
		}
	}

	return nil
}

func writeEmptyArrayToFile(filePath string) error {
	return writeToFile(filePath, []interface{}{})
}

func writeToFile(filePath string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, file, 0o644)
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
