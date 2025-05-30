package repository

import (
	"database/sql"
	"errors"
	"frappuccino/models"

	"github.com/lib/pq"
)

var ErrDuplicateID = errors.New("duplicate ID")

type MenuRepository interface {
	Create(item models.MenuItem) (models.MenuItem, error)
	GetAll() ([]models.MenuItem, error)
	GetByID(id int64) (models.MenuItem, error)
	Update(id int64, item models.MenuItem) (models.MenuItem, error)
	Delete(id int64) error
}

type menuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(item models.MenuItem) (models.MenuItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return item, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `INSERT INTO menu_items (product_name, description, categories, price) 
	          VALUES ($1, $2, $3, $4) RETURNING product_id`
	err = tx.QueryRow(query, item.Name, item.Description, pq.Array(item.Categories), item.Price).Scan(&item.ID)
	if err != nil {
		return item, err
	}

	for _, ing := range item.Ingredients {
		ingQuery := `INSERT INTO menu_item_ingredients (ingredient_id, product_id, quantity) VALUES ($1, $2, $3)`
		_, err = tx.Exec(ingQuery, ing.IngredientID, item.ID, ing.Quantity)
		if err != nil {
			return item, err
		}
	}

	return item, nil
}

func (r *menuRepository) GetAll() ([]models.MenuItem, error) {
	query := `SELECT product_id, product_name, description, categories, price FROM menu_items`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, pq.Array(&item.Categories), &item.Price); err != nil {
			return nil, err
		}

		ingQuery := `SELECT i.ingredient_id, i.name, mi.quantity
		             FROM menu_item_ingredients mi
		             JOIN inventory i ON i.ingredient_id = mi.ingredient_id
		             WHERE mi.product_id = $1`

		ingRows, err := r.db.Query(ingQuery, item.ID)
		if err != nil {
			return nil, err
		}

		var ingredients []models.MenuItemIngredient
		for ingRows.Next() {
			var ing models.MenuItemIngredient
			if err := ingRows.Scan(&ing.IngredientID, &ing.ProductName, &ing.Quantity); err != nil {
				ingRows.Close()
				return nil, err
			}
			ingredients = append(ingredients, ing)
		}
		ingRows.Close()

		item.Ingredients = ingredients
		items = append(items, item)
	}

	return items, nil
}

func (r *menuRepository) GetByID(id int64) (models.MenuItem, error) {
	query := `SELECT product_id, product_name, description, categories, price FROM menu_items WHERE product_id = $1`
	var item models.MenuItem
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Description, pq.Array(&item.Categories), &item.Price)
	if err == sql.ErrNoRows {
		return models.MenuItem{}, ErrNotFound
	}
	if err != nil {
		return models.MenuItem{}, err
	}

	ingQuery := `SELECT i.ingredient_id, i.name, mi.quantity
				 FROM menu_item_ingredients mi
				 JOIN inventory i ON i.ingredient_id = mi.ingredient_id
				 WHERE mi.product_id = $1`

	rows, err := r.db.Query(ingQuery, item.ID)
	if err != nil {
		return item, err
	}
	defer rows.Close()

	for rows.Next() {
		var ing models.MenuItemIngredient
		if err := rows.Scan(&ing.IngredientID, &ing.ProductName, &ing.Quantity); err != nil {
			return item, err
		}
		item.Ingredients = append(item.Ingredients, ing)
	}

	return item, nil
}

func (r *menuRepository) Update(id int64, item models.MenuItem) (models.MenuItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.MenuItem{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `UPDATE menu_items SET product_name = $1, description = $2, categories = $3, price = $4, updated_at = NOW() WHERE product_id = $5`
	result, err := tx.Exec(query, item.Name, item.Description, pq.Array(item.Categories), item.Price, id)
	if err != nil {
		return models.MenuItem{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return models.MenuItem{}, ErrNotFound
	}

	_, err = tx.Exec(`DELETE FROM menu_item_ingredients WHERE product_id = $1`, id)
	if err != nil {
		return models.MenuItem{}, err
	}

	for _, ing := range item.Ingredients {
		ingQuery := `INSERT INTO menu_item_ingredients (ingredient_id, product_id, quantity) VALUES ($1, $2, $3)`
		_, err = tx.Exec(ingQuery, ing.IngredientID, id, ing.Quantity)
		if err != nil {
			return models.MenuItem{}, err
		}
	}

	item.ID = id
	return item, nil
}

func (r *menuRepository) Delete(id int64) error {
	query := `DELETE FROM menu_items WHERE product_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
