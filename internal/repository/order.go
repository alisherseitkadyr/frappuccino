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
	GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error)
	CreateTx(tx *sql.Tx, order models.Order) (models.Order, error)
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
	tx, err := r.db.Begin()
	if err != nil {
		return models.Order{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Вставляем основной заказ
	query := `
		INSERT INTO orders (customer_name, total_price)
		VALUES ($1, $2)
		RETURNING order_id, created_at`
	err = tx.QueryRow(query, order.CustomerName, order.TotalPrice).
		Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return models.Order{}, err
	}

	// Вставляем товары заказа
	for _, item := range order.Items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			return models.Order{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) GetAll() ([]models.Order, error) {
	query := `SELECT order_id, customer_name, total_price, created_at FROM orders ORDER BY order_id DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.TotalPrice, &order.CreatedAt); err != nil {
			return nil, err
		}

		order.Items, err = r.getOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}

func (r *orderRepository) GetByID(id int64) (models.Order, error) {
	var order models.Order

	query := `SELECT order_id, customer_name, total_price, created_at FROM orders WHERE order_id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.CustomerName, &order.TotalPrice, &order.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Order{}, ErrNotFound
	}
	if err != nil {
		return models.Order{}, err
	}

	order.Items, err = r.getOrderItems(order.ID)
	if err != nil {
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
	result, err := r.db.Exec(query, updatedOrder.CustomerName, itemsJSON, updatedOrder.TotalPrice, id)
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
	// 1. Вставка заказа
	query := `INSERT INTO orders (customer_name, status, created_at) VALUES ($1, $2, $3) RETURNING order_id`
	err := tx.QueryRow(query, order.CustomerName, order.Status, order.CreatedAt).Scan(&order.ID)
	if err != nil {
		return models.Order{}, err
	}

	// 2. Подсчёт необходимых ингредиентов
	ingredientNeeds := make(map[int64]int) // ingredient_id -> total needed
	for _, item := range order.Items {
		ingredients, err := r.GetIngredientsByProductID(item.ProductID)
		if err != nil {
			return models.Order{}, err
		}
		for _, ing := range ingredients {
			ingredientNeeds[ing.IngredientID] += ing.Quantity * item.Quantity
		}

		// вставка в order_items
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`
		_, err = tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			return models.Order{}, err
		}
	}

	// 3. Проверка и обновление inventory
	for ingID, totalQty := range ingredientNeeds {
		err := r.UpdateInventory(tx, ingID, totalQty)
		if err != nil {
			return models.Order{}, err
		}
	}

	return order, nil
}

func (r *orderRepository) GetIngredientsByProductID(productID int64) ([]models.MenuItemIngredient, error) {
	query := `
	SELECT mi.ingredient_id, inv.name, mi.quantity
	FROM menu_item_ingredients mi
	JOIN inventory inv ON inv.ingredient_id = mi.ingredient_id
	WHERE mi.product_id = $1`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.MenuItemIngredient
	for rows.Next() {
		var ing models.MenuItemIngredient
		if err := rows.Scan(&ing.IngredientID, &ing.ProductName, &ing.Quantity); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients, nil
}

func (r *orderRepository) UpdateInventory(tx *sql.Tx, ingredientID int64, quantity int) error {
	query := `UPDATE inventory SET quantity = quantity - $1 WHERE ingredient_id = $2 AND quantity >= $1`
	result, err := tx.Exec(query, quantity, ingredientID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("not enough inventory")
	}
	return nil
}

func (r *orderRepository) getOrderItems(orderID int64) ([]models.OrderItem, error) {
	query := `
		SELECT oi.product_id, p.product_name, oi.quantity
		FROM order_items oi
		JOIN menu_items p ON oi.product_id = p.product_id
		WHERE oi.order_id = $1`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *orderRepository) GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error) {
	query := `
	SELECT mi.product_name, SUM(oi.quantity) as total_quantity
	FROM order_items oi
	JOIN menu_items mi ON oi.product_id = mi.product_id
	JOIN orders o ON oi.order_id = o.order_id
	WHERE ($1 = '' OR o.created_at >= $1::timestamptz)
	AND ($2 = '' OR o.created_at <= $2::timestamptz)
	GROUP BY mi.product_name;
`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var name string
		var quantity int
		if err := rows.Scan(&name, &quantity); err != nil {
			return nil, err
		}
		result[name] = quantity
	}

	return result, nil
}
