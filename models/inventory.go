// models/inventory.go
package models

import "time"

type InventoryItem struct {
	IngredientID string  `json:"ingredient_id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

type InventoryItemRequest struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

func NewInventoryItem(name string, quantity float64, unit string) InventoryItem {
	return InventoryItem{
		IngredientID: generateIngredientID(),
		Name:         name,
		Quantity:     quantity,
		Unit:         unit,
	}
}

func generateIngredientID() string {
	return "ingr_" + time.Now().Format("20060102150405")
}
