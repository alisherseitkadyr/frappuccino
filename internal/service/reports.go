package service

import (
	"frappuccino/internal/repository"
	"frappuccino/models"
	"sort"
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

// GetTotalSales подсчитывает сумму выручки по всем закрытым заказам
func (s *reportsService) GetTotalSales() (float64, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return 0, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return 0, err
	}

	// Создаем мапу для быстрого поиска цены по ID товара
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

// GetPopularItems возвращает топ-N популярных товаров по количеству заказов
func (s *reportsService) GetPopularItems(limit int) ([]models.MenuItem, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Подсчитываем количество заказов каждого товара
	popularity := make(map[int64]int)
	for _, order := range orders {
		if order.Status == "closed" {
			for _, item := range order.Items {
				popularity[item.ProductID] += item.Quantity
			}
		}
	}

	// Создаем срез с данными о популярности для сортировки
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

	// Сортируем по убыванию популярности
	sort.Slice(popularItems, func(i, j int) bool {
		return popularItems[i].OrderCount > popularItems[j].OrderCount
	})

	// Формируем результат с ограничением по limit
	if limit > len(popularItems) {
		limit = len(popularItems)
	}
	result := make([]models.MenuItem, limit)
	for i := 0; i < limit; i++ {
		result[i] = popularItems[i].MenuItem
	}

	return result, nil
}
