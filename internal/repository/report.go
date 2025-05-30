package repository

import (
	"database/sql"
	"fmt"
	"frappuccino/models"
	"strings"
	"time"
)

type ReportRepository interface {
	SearchReports(query string, filters []string, minPrice, maxPrice float64) (models.SearchReportResponse, error)
	GetOrderedItemsByDay(year int, month time.Month) (map[int]int, error)
	GetOrderedItemsByMonth(year int) (map[string]int, error)
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) SearchReports(query string, filters []string, minPrice, maxPrice float64) (models.SearchReportResponse, error) {
	response := models.SearchReportResponse{}
	searchQuery := fmt.Sprintf("plainto_tsquery('english', $1)")

	if contains(filters, "menu") || contains(filters, "all") {
		menuSQL := `
			SELECT product_id, product_name, description, price,
			ts_rank_cd(to_tsvector('english', product_name || ' ' || description), ` + searchQuery + `) AS relevance
			FROM menu_items
			WHERE to_tsvector('english', product_name || ' ' || description) @@ ` + searchQuery + `
			AND price BETWEEN $2 AND $3
			ORDER BY relevance DESC;
		`

		rows, err := r.db.Query(menuSQL, query, minPrice, maxPrice)
		if err != nil {
			return response, err
		}
		defer rows.Close()

		for rows.Next() {
			var item models.MenuItemSearchResult
			if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Relevance); err != nil {
				return response, err
			}
			response.MenuItems = append(response.MenuItems, item)
		}
	}

	if contains(filters, "orders") || contains(filters, "all") {
		orderSQL := `
			SELECT o.order_id, o.customer_name, o.total_price,
			ts_rank_cd(to_tsvector('english', o.customer_name), ` + searchQuery + `) AS relevance
			FROM orders o
			WHERE to_tsvector('english', o.customer_name) @@ ` + searchQuery + `
			AND o.total_price BETWEEN $2 AND $3
			ORDER BY relevance DESC;
		`

		rows, err := r.db.Query(orderSQL, query, minPrice, maxPrice)
		if err != nil {
			return response, err
		}
		defer rows.Close()

		for rows.Next() {
			var ord models.OrderSearchResult
			if err := rows.Scan(&ord.ID, &ord.CustomerName, &ord.Total, &ord.Relevance); err != nil {
				return response, err
			}
			// fetch item names for order
			itemsRows, err := r.db.Query(`
				SELECT mi.product_name
				FROM order_items oi
				JOIN menu_items mi ON oi.product_id = mi.product_id
				WHERE oi.order_id = $1
			`, ord.ID)
			if err != nil {
				return response, err
			}
			for itemsRows.Next() {
				var name string
				itemsRows.Scan(&name)
				ord.Items = append(ord.Items, name)
			}
			itemsRows.Close()
			response.Orders = append(response.Orders, ord)
		}
	}

	response.TotalMatches = len(response.MenuItems) + len(response.Orders)
	return response, nil
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if strings.TrimSpace(v) == value {
			return true
		}
	}
	return false
}

func (r *reportRepository) GetOrderedItemsByDay(year int, month time.Month) (map[int]int, error) {
	query := `
		SELECT EXTRACT(DAY FROM day) AS day, COUNT(*) AS count
		FROM (
  		SELECT DATE(created_at) AS day
  		FROM orders
  		WHERE created_at >= DATE_TRUNC('month', MAKE_DATE($1, $2, 1))
    	AND created_at <  DATE_TRUNC('month', MAKE_DATE($1, $2, 1)) + INTERVAL '1 month') sub
		GROUP BY day
		ORDER BY day;`
	rows, err := r.db.Query(query, year, int(month))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var day int
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result[day] = count
	}
	return result, nil
}

func (r *reportRepository) GetOrderedItemsByMonth(year int) (map[string]int, error) {
	query := `
		SELECT TO_CHAR(created_at, 'Month') as month, COUNT(*) as count
		FROM orders
		WHERE EXTRACT(YEAR FROM created_at) = $1
		GROUP BY month
		ORDER BY MIN(created_at)`
	rows, err := r.db.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var month string
		var count int
		if err := rows.Scan(&month, &count); err != nil {
			return nil, err
		}
		month = strings.TrimSpace(strings.ToLower(month)) // например "january"
		result[month] = count
	}
	return result, nil
}
