package models

import "time"

type InventoryItem struct {
    IngredientID int64  
    Name         string  
    Quantity     float64
    Unit         string  
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type InventoryItemRequest struct {
	Name     string  
	Quantity float64 
	Unit     string  
}

// NewInventoryItem создаёт новый InventoryItem с уникальным int64 ID
func NewInventoryItem(name string, quantity float64, unit string) InventoryItem {
    return InventoryItem{
        IngredientID: 0, // заполняется после вставки в БД
        Name:         name,
        Quantity:     quantity,
        Unit:         unit,
    }
}

