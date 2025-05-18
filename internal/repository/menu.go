package repository

import (
	"database/sql"
	"errors"
	"frappuccino/models"
)

// ErrNotFound    = errors.New("not found")
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
	query := `INSERT INTO menu_items (name, description, price) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, item.Name, item.Description, item.Price).Scan(&item.ID)
	return item, err
}

func (r *menuRepository) GetAll() ([]models.MenuItem, error) {
	query := `SELECT id, name, description, price FROM menu_items`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *menuRepository) GetByID(id int64) (models.MenuItem, error) {
	query := `SELECT id, name, description, price FROM menu_items WHERE id = $1`
	var item models.MenuItem
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Description, &item.Price)
	if err == sql.ErrNoRows {
		return models.MenuItem{}, ErrNotFound
	}
	return item, err
}

func (r *menuRepository) Update(id int64, item models.MenuItem) (models.MenuItem, error) {
	query := `UPDATE menu_items SET name = $1, description = $2, price = $3 WHERE id = $4`
	result, err := r.db.Exec(query, item.Name, item.Description, item.Price, id)
	if err != nil {
		return models.MenuItem{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return models.MenuItem{}, ErrNotFound
	}
	item.ID = id
	return item, nil
}

func (r *menuRepository) Delete(id int64) error {
	query := `DELETE FROM menu_items WHERE id = $1`
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
