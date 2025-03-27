package service

import (
	"errors"
	"fmt"
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
	// Validate required fields
	if order.CustomerName == "" {
		return models.Order{}, errors.New("customer_name is required")
	}

	if len(order.Items) == 0 {
		return models.Order{}, errors.New("order must contain at least one item")
	}

	// Validate each item
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return models.Order{}, errors.New("quantity must be positive")
		}

		menuItem, err := s.menuRepo.GetByID(item.ProductID)
		if err != nil {
			slog.Error("Invalid product ID",
				"product_id", item.ProductID,
				"error", err)
			return models.Order{}, fmt.Errorf("product ID '%s' not found in menu", item.ProductID)
		}

		// Check ingredients for each menu item
		for _, ingredient := range menuItem.Ingredients {
			invItem, err := s.inventoryRepo.GetByID(ingredient.IngredientID)
			if err != nil {
				slog.Error("Inventory item not found",
					"ingredient_id", ingredient.IngredientID,
					"error", err)
				return models.Order{}, fmt.Errorf("ingredient '%s' not available", ingredient.IngredientID)
			}

			needed := ingredient.Quantity * float64(item.Quantity)
			if invItem.Quantity < needed {
				return models.Order{}, fmt.Errorf(
					"not enough %s. Need %.2f%s, have %.2f%s",
					invItem.Name,
					needed,
					invItem.Unit,
					invItem.Quantity,
					invItem.Unit,
				)
			}
		}
	}

	// Generate a new order ID
	order.ID = generateOrderID()
	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Deduct inventory (only if all checks pass)
	for _, item := range order.Items {
		menuItem, _ := s.menuRepo.GetByID(item.ProductID)
		for _, ingredient := range menuItem.Ingredients {
			invItem, _ := s.inventoryRepo.GetByID(ingredient.IngredientID)
			invItem.Quantity -= ingredient.Quantity * float64(item.Quantity)
			if _, err := s.inventoryRepo.Update(invItem); err != nil {
				slog.Error("Failed to update inventory", "error", err)
				return models.Order{}, fmt.Errorf("failed to update inventory: %v", err)
			}
		}
	}

	// Create order
	createdOrder, err := s.orderRepo.Create(order)
	if err != nil {
		slog.Error("Failed to save order", "error", err)
		return models.Order{}, errors.New("failed to save order")
	}

	return createdOrder, nil
}

// generateOrderID creates a unique order ID using current timestamp
func generateOrderID() string {
	return fmt.Sprintf("order_%d", time.Now().UnixNano())
}

func formatQuantity(quantity float64, unit string) string {
	if quantity == float64(int(quantity)) {
		return fmt.Sprintf("%d%s", int(quantity), unit)
	}
	return fmt.Sprintf("%.2f%s", quantity, unit)
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

	if order.Status == "closed" {
		return models.Order{}, errors.New("order is already closed")
	}

	order.Status = "closed"
	return s.orderRepo.Update(id, order)
}
