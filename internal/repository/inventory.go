package repository

import (
	"database/sql"
	"fmt"
	"frappuccino/models"
	"time"
)

type InventoryRepository interface {
	Create(item models.InventoryItem) (models.InventoryItem, error)
	GetAll() ([]models.InventoryItem, error)
	GetByID(id int64) (models.InventoryItem, error)
	Update(item models.InventoryItem) (models.InventoryItem, error)
	UpdateTx(tx *sql.Tx, item models.InventoryItem) (models.InventoryItem, error)
	Delete(id int64) error
	GetLeftOvers(sortBy string, offset, limit int) ([]models.InventoryItem, int, error)
}

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(item models.InventoryItem) (models.InventoryItem, error) {
	query := `INSERT INTO inventory (name, quantity, unit, created_at) VALUES ($1, $2, $3, $4) RETURNING ingredient_id`
	createdAt := time.Now()
	err := r.db.QueryRow(query, item.Name, item.Quantity, item.Unit, createdAt).Scan(&item.IngredientID)
	item.CreatedAt=createdAt
	return item, err
}

func (r *inventoryRepository) GetAll() ([]models.InventoryItem, error) {
	query := `SELECT ingredient_id, name, quantity, unit, created_at, updated_at FROM inventory`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *inventoryRepository) GetByID(id int64) (models.InventoryItem, error) {
	query := `SELECT ingredient_id, name, quantity, unit, created_at, updated_at FROM inventory WHERE ingredient_id = $1`
	var item models.InventoryItem
	err := r.db.QueryRow(query, id).Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit,  &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return models.InventoryItem{}, ErrNotFound
	}
	return item, err
}

func (r *inventoryRepository) Update(item models.InventoryItem) (models.InventoryItem, error) {
	query := `UPDATE inventory SET name = $1, quantity = $2, unit = $3, updated_at = $5 WHERE ingredient_id = $4`
	updated_at:=time.Now()
	result, err := r.db.Exec(query, item.Name, item.Quantity, item.Unit, item.IngredientID, updated_at)
	item.UpdatedAt=updated_at
	if err != nil {
		return models.InventoryItem{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return models.InventoryItem{}, ErrNotFound
	}
	return item, nil
}

func (r *inventoryRepository) Delete(id int64) error {
	query := `DELETE FROM inventory WHERE ingredient_id = $1`
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

func (r *inventoryRepository) UpdateTx(tx *sql.Tx, item models.InventoryItem) (models.InventoryItem, error) {
	query := "UPDATE inventory SET quantity = $1, unit = $2, updated_at = CURRENT_TIMESTAMP WHERE ingredient_id = $3"
	_, err := tx.Exec(query, item.Quantity, item.Unit, item.IngredientID)
	if err != nil {
		return models.InventoryItem{}, err
	}
	return item, nil
}

func (r *inventoryRepository) GetLeftOvers(sortBy string, offset, limit int) ([]models.InventoryItem, int, error) {
	validSortFields := map[string]string{
		"price":    "price",
		"quantity": "quantity",
	}

	sortField, ok := validSortFields[sortBy]
	if !ok {
		sortField = "ingredient_id" // default sort
	}

	query := fmt.Sprintf(`SELECT ingredient_id, name, quantity, unit FROM inventory ORDER BY %s DESC LIMIT ? OFFSET ?`, sortField)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	// Подсчёт общего количества для пагинации
	var total int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM inventory`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
