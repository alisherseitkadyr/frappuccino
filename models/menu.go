package models

import (
	"time"
)

type MenuItem struct {
	ID          int64
	Name        string
	Description string
	Categories  []string
	Price       float64
	Ingredients []MenuItemIngredient
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MenuItemIngredient struct {
	IngredientID int64
	ProductName  string
	Quantity     int
}

// type MenuItemRequest struct {
// 	Name        string
// 	Description string
// 	Categories  []string
// 	Price       float64
// 	Ingredients []MenuItemIngredient
// }

func NewMenuItem(name, description string, categories []string, price float64, ingredients []MenuItemIngredient) MenuItem {
	return MenuItem{
		ID:          0,
		Name:        name,
		Description: description,
		Categories:  categories,
		Price:       price,
		Ingredients: ingredients,
	}
}
