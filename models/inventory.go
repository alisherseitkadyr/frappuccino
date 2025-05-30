package models

import "time"

type InventoryItem struct {
	IngredientID int64
	Name         string
	Quantity     int
	Unit         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewInventoryItem(name string, quantity int, unit string) InventoryItem {
	return InventoryItem{
		IngredientID: 0, // заполняется после вставки в БД
		Name:         name,
		Quantity:     quantity,
		Unit:         unit,
	}
}
