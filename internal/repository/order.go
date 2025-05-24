package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"frappuccino/models"
)

var ErrNotFound = errors.New("not found")

type OrderRepository interface {
	Create(order models.Order) (models.Order, error)
	CreateTx(tx *sql.Tx, order models.Order) (models.Order, error) // 👈 добавь этот метод
	GetAll() ([]models.Order, error)
	GetByID(id int64) (models.Order, error)
	Update(id int64, order models.Order) (models.Order, error)
	Delete(id int64) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order models.Order) (models.Order, error) {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return models.Order{}, err
	}

	query := `
        INSERT INTO orders (customer_name, items, total_price)
        VALUES ($1, $2, $3)
        RETURNING order_id, created_at`
	err = r.db.QueryRow(query, order.Customer_name, itemsJSON, order.TotalPrice).
		Scan(&order.ID, &order.CreatedAt)
		
	return order, err
}

func (r *orderRepository) GetAll() ([]models.Order, error) {
	query := `SELECT order_id, customer_name, items, total_price, created_at FROM orders ORDER BY order_id DESC`
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var itemsData []byte

		if err := rows.Scan(&order.ID, &order.Customer_name, &itemsData, &order.TotalPrice, &order.CreatedAt); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(itemsData, &order.Items); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}

func (r *orderRepository) GetByID(id int64) (models.Order, error) {
	var order models.Order
	var itemsData []byte

	query := `SELECT order_id, customer_name, items, total_price, created_at FROM orders WHERE order_id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.Customer_name, &itemsData, &order.TotalPrice, &order.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Order{}, ErrNotFound
	}
	if err != nil {
		return models.Order{}, err
	}

	if err := json.Unmarshal(itemsData, &order.Items); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) Update(id int64, updatedOrder models.Order) (models.Order, error) {
	itemsJSON, err := json.Marshal(updatedOrder.Items)
	if err != nil {
		return models.Order{}, err
	}

	query := `
        UPDATE orders
        SET customer_name = $1, items = $2, total_price = $3
        WHERE order_id = $4`
	result, err := r.db.Exec(query, updatedOrder.Customer_name, itemsJSON, updatedOrder.TotalPrice, id)
	if err != nil {
		return models.Order{}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return models.Order{}, ErrNotFound
	}

	updatedOrder.ID = id
	return updatedOrder, nil
}

func (r *orderRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM orders WHERE order_id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *orderRepository) CreateTx(tx *sql.Tx, order models.Order) (models.Order, error) {
	query := "INSERT INTO orders (customer_name, status, created_at) VALUES (?, ?, ?) RETURNING order_id"
	err := tx.QueryRow(query, order.Customer_name, order.Status, order.CreatedAt).Scan(&order.ID)
	if err != nil {
		return models.Order{}, err
	}

	for _, item := range order.Items {
		itemQuery := "INSERT INTO order_items (order_id, product_id, quantity) VALUES (?, ?, ?)"
		_, err := tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			return models.Order{}, err
		}
	}
	return order, nil
}
