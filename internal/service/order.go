package service

import (
	"errors"
	"hot-coffee/internal/repository"
	"hot-coffee/models"
	"log/slog"
	"time"
)

type OrderService interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrders() ([]models.Order, error)
	GetOrder(id string) (models.Order, error)
	UpdateOrder(id string, order models.Order) (models.Order, error)
	DeleteOrder(id string) error
	CloseOrder(id string) (models.Order, error)
}

type orderService struct {
	orderRepo     repository.OrderRepository
	menuRepo      repository.MenuRepository
	inventoryRepo repository.InventoryRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	menuRepo repository.MenuRepository,
	inventoryRepo repository.InventoryRepository,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *orderService) CreateOrder(order models.Order) (models.Order, error) {
	// Validate order items
	for _, item := range order.Items {
		_, err := s.menuRepo.GetByID(item.ProductID)
		if err != nil {
			slog.Error("Invalid product ID in order", "product_id", item.ProductID, "error", err)
			return models.Order{}, errors.New("invalid product ID in order items")
		}
	}

	// Check inventory
	ingredientsNeeded := make(map[string]float64)
	for _, item := range order.Items {
		menuItem, err := s.menuRepo.GetByID(item.ProductID)
		if err != nil {
			return models.Order{}, err
		}

		for _, ingredient := range menuItem.Ingredients {
			ingredientsNeeded[ingredient.IngredientID] += ingredient.Quantity * float64(item.Quantity)
		}
	}

	// Verify inventory levels
	for ingredientID, needed := range ingredientsNeeded {
		inventoryItem, err := s.inventoryRepo.GetByID(ingredientID)
		if err != nil {
			return models.Order{}, err
		}

		if inventoryItem.Quantity < needed {
			return models.Order{}, errors.New("insufficient inventory for ingredient " + inventoryItem.Name)
		}
	}

	// Deduct inventory
	for ingredientID, needed := range ingredientsNeeded {
		inventoryItem, err := s.inventoryRepo.GetByID(ingredientID)
		if err != nil {
			return models.Order{}, err
		}

		inventoryItem.Quantity -= needed
		if _, err := s.inventoryRepo.Update(inventoryItem); err != nil {
			return models.Order{}, err
		}
	}

	// Set order defaults
	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Create order
	return s.orderRepo.Create(order)
}

func (s *orderService) GetOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *orderService) GetOrder(id string) (models.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *orderService) UpdateOrder(id string, order models.Order) (models.Order, error) {
	existingOrder, err := s.orderRepo.GetByID(id)
	if err != nil {
		return models.Order{}, err
	}

	// Preserve some fields
	order.ID = existingOrder.ID
	order.CreatedAt = existingOrder.CreatedAt

	return s.orderRepo.Update(id, order)
}

func (s *orderService) DeleteOrder(id string) error {
	return s.orderRepo.Delete(id)
}

func (s *orderService) CloseOrder(id string) (models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return models.Order{}, err
	}

	order.Status = "closed"
	return s.orderRepo.Update(id, order)
}
