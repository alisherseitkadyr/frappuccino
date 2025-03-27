// models/menu.go
package models

import "time"

type MenuItem struct {
	ID          string               `json:"product_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
}

type MenuItemIngredient struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}

type MenuItemRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
}

func NewMenuItem(name, description string, price float64, ingredients []MenuItemIngredient) MenuItem {
	return MenuItem{
		ID:          generateProductID(),
		Name:        name,
		Description: description,
		Price:       price,
		Ingredients: ingredients,
	}
}

func generateProductID() string {
	return "prod_" + time.Now().Format("20060102150405")
}
