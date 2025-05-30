package service

import (
	"fmt"
	"frappuccino/internal/repository"
	"frappuccino/models"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ReportsService interface {
	SearchReport(q, filter, min, max string) (models.SearchReportResponse, error)
	GetOrderedItemsByPeriod(period, month, year string) (models.OrderedItemsByPeriodResponse, error)
	GetTotalSales() (float64, error)
	GetPopularItems(limit int) ([]models.MenuItem, error)
}

type reportsService struct {
	orderRepo repository.OrderRepository
	menuRepo  repository.MenuRepository
	repo      repository.ReportRepository
}

func NewReportsService(
	orderRepo repository.OrderRepository,
	menuRepo repository.MenuRepository,
	reportRepo repository.ReportRepository,
) ReportsService {
	return &reportsService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
		repo:      reportRepo,
	}
}

func (s *reportsService) GetTotalSales() (float64, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return 0, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return 0, err
	}

	menuItemMap := make(map[int64]models.MenuItem, len(menuItems))
	for _, item := range menuItems {
		menuItemMap[item.ID] = item
	}

	var totalSales float64
	for _, order := range orders {
		if order.Status == "closed" {
			for _, item := range order.Items {
				if menuItem, exists := menuItemMap[item.ProductID]; exists {
					totalSales += menuItem.Price * float64(item.Quantity)
				}
			}
		}
	}

	return totalSales, nil
}

func (s *reportsService) GetPopularItems(limit int) ([]models.MenuItem, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}

	popularity := make(map[int64]int)
	for _, order := range orders {
		if order.Status == "closed" {
			for _, item := range order.Items {
				popularity[item.ProductID] += item.Quantity
			}
		}
	}

	type popularItem struct {
		MenuItem   models.MenuItem
		OrderCount int
	}

	popularItems := make([]popularItem, 0, len(menuItems))
	for _, menuItem := range menuItems {
		count := popularity[menuItem.ID]
		popularItems = append(popularItems, popularItem{
			MenuItem:   menuItem,
			OrderCount: count,
		})
	}

	sort.Slice(popularItems, func(i, j int) bool {
		return popularItems[i].OrderCount > popularItems[j].OrderCount
	})

	if limit > len(popularItems) {
		limit = len(popularItems)
	}
	result := make([]models.MenuItem, limit)
	for i := 0; i < limit; i++ {
		result[i] = popularItems[i].MenuItem
	}

	return result, nil
}

func (s *reportsService) SearchReport(q, filter, min, max string) (models.SearchReportResponse, error) {
	minPrice := 0.0
	maxPrice := 999999.0
	if s.repo == nil {
		log.Fatal("s.repo is nil")
	}
	if min != "" {
		if p, err := strconv.ParseFloat(min, 64); err == nil {
			minPrice = p
		}
	}
	if max != "" {
		if p, err := strconv.ParseFloat(max, 64); err == nil {
			maxPrice = p
		}
	}

	filters := []string{"all"}
	if filter != "" {
		filters = strings.Split(filter, ",")
	}

	return s.repo.SearchReports(q, filters, minPrice, maxPrice)
}

func (s *reportsService) GetOrderedItemsByPeriod(period, month, year string) (models.OrderedItemsByPeriodResponse, error) {
	resp := models.OrderedItemsByPeriodResponse{
		Period: period,
	}

	if period == "day" {
		if month == "" {
			return resp, fmt.Errorf("month parameter required for period=day")
		}
		monthParsed, err := time.Parse("January", strings.Title(month))
		if err != nil {
			return resp, fmt.Errorf("invalid month: %s", month)
		}

		yearInt := time.Now().Year()
		if year != "" {
			if y, err := strconv.Atoi(year); err == nil {
				yearInt = y
			}
		}

		resp.Month = strings.ToLower(month)

		data, err := s.repo.GetOrderedItemsByDay(yearInt, monthParsed.Month())
		if err != nil {
			return resp, err
		}

		// Формируем массив по дням месяца
		var items []models.OrderedItemCount
		daysInMonth := 31 // можно оптимизировать, но для простоты 31
		for i := 1; i <= daysInMonth; i++ {
			count := data[i]
			items = append(items, models.OrderedItemCount{
				Key:   fmt.Sprintf("%d", i),
				Count: count,
			})
		}
		resp.OrderedItems = items
		return resp, nil
	}

	if period == "month" {
		yearInt := time.Now().Year()
		if year != "" {
			if y, err := strconv.Atoi(year); err == nil {
				yearInt = y
			}
		}

		resp.Year = year

		data, err := s.repo.GetOrderedItemsByMonth(yearInt)
		if err != nil {
			return resp, err
		}

		// Формируем массив месяцев в нужном порядке
		monthsOrder := []string{
			"january", "february", "march", "april", "may", "june",
			"july", "august", "september", "october", "november", "december",
		}

		var items []models.OrderedItemCount
		for _, m := range monthsOrder {
			count := data[m]
			items = append(items, models.OrderedItemCount{
				Key:   m,
				Count: count,
			})
		}
		resp.OrderedItems = items
		return resp, nil
	}

	return resp, fmt.Errorf("invalid period parameter: %s", period)
}
