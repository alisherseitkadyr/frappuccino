package models

type MenuItemSearchResult struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Relevance   float64 `json:"relevance"`
}

type OrderSearchResult struct {
	ID           int      `json:"id"`
	CustomerName string   `json:"customer_name"`
	Items        []string `json:"items"`
	Total        float64  `json:"total"`
	Relevance    float64  `json:"relevance"`
}

type SearchReportResponse struct {
	MenuItems     []MenuItemSearchResult `json:"menu_items"`
	Orders        []OrderSearchResult    `json:"orders"`
	TotalMatches  int                    `json:"total_matches"`
}


type OrderedItemCount struct {
	Key   string `json:"key"` 
	Count int    `json:"count"` 
}

type OrderedItemsByPeriodResponse struct {
	Period       string             `json:"period"`
	Month        string             `json:"month,omitempty"`
	Year         string             `json:"year,omitempty"`
	OrderedItems []OrderedItemCount `json:"orderedItems"`
}