package service

import (
	"hot-coffee/models"
	"hot-coffee/internal/repository"
)

type ReportsService interface {
	GetTotalSales() (float64, error)
	GetPopularItems(limit int) ([]models.MenuItem, error)
}

type reportsService struct {
	orderRepo repository.OrderRepository
	menuRepo  repository.MenuRepository
}

func NewReportsService(
	orderRepo repository.OrderRepository,
	menuRepo repository.MenuRepository,
) ReportsService {
	return &reportsService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
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

	menuItemMap := make(map[string]models.MenuItem)
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

	// Count how many times each item was ordered
	popularity := make(map[string]int)
	for _, order := range orders {
		if order.Status == "closed" {
			for _, item := range order.Items {
				popularity[item.ProductID] += item.Quantity
			}
		}
	}

	// Create a list of menu items with their popularity
	var popularItems []struct {
		MenuItem  models.MenuItem
		OrderCount int
	}

	for _, menuItem := range menuItems {
		popularItems = append(popularItems, struct {
			MenuItem  models.MenuItem
			OrderCount int
		}{
			MenuItem:  menuItem,
			OrderCount: popularity[menuItem.ID],
		})
	}

	// Sort by popularity (simplified for this example)
	// In a real implementation, we'd use a proper sorting algorithm
	for i := 0; i < len(popularItems); i++ {
		for j := i + 1; j < len(popularItems); j++ {
			if popularItems[i].OrderCount < popularItems[j].OrderCount {
				popularItems[i], popularItems[j] = popularItems[j], popularItems[i]
			}
		}
	}

	// Return top N items
	result := make([]models.MenuItem, 0, limit)
	for i := 0; i < limit && i < len(popularItems); i++ {
		result = append(result, popularItems[i].MenuItem)
	}

	return result, nil
}