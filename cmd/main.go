package main

import (
	"database/sql"
	"flag"
	"fmt"
	"frappuccino/internal/api"
	"frappuccino/internal/repository"
	"frappuccino/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	port := flag.Int("port", 8090, "Port number")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		printUsage()
		return
	}

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to open database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Verify the connection is alive
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database", "error", err)
		os.Exit(1)
	}

	// Initialize repositories using the DB connection
	orderRepo := repository.NewOrderRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Initialize services
	orderSvc := service.NewOrderService(orderRepo, menuRepo, inventoryRepo, db)
	menuSvc := service.NewMenuService(menuRepo)
	inventorySvc := service.NewInventoryService(inventoryRepo)
	reportsSvc := service.NewReportsService(orderRepo, menuRepo, reportRepo)

	// Initialize router
	router := api.NewRouter(orderSvc, menuSvc, inventorySvc, reportsSvc)

	log.Printf("Starting server", "port", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		log.Printf("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Coffee Shop Management System

Usage:
  frappuccino [--port <N>] [--db <connection-string>]
  frappuccino --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --db S       PostgreSQL connection string.`)
}
