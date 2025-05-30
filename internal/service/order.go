package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"frappuccino/internal/repository"
	"frappuccino/models"
	"log"
)

type OrderService interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrders() ([]models.Order, error)
	GetOrder(id int64) (models.Order, error)
	GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error)
	UpdateOrder(id int64, order models.Order) (models.Order, error)
	DeleteOrder(id int64) error
	CloseOrder(id int64) (models.Order, error)
}

type orderService struct {
	orderRepo     repository.OrderRepository
	menuRepo      repository.MenuRepository
	inventoryRepo repository.InventoryRepository
	db            *sql.DB
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	menuRepo repository.MenuRepository,
	inventoryRepo repository.InventoryRepository,
	db *sql.DB,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
		db:            db,
	}
}

func (s *orderService) CreateOrder(order models.Order) (models.Order, error) {
	if order.CustomerName == "" {
		return models.Order{}, errors.New("customer_name is required")
	}
	if len(order.Items) == 0 {
		return models.Order{}, errors.New("order must contain at least one item")
	}

	menuCache := make(map[int64]models.MenuItem)

	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return models.Order{}, errors.New("quantity must be positive")
		}

		menuItem, ok := menuCache[item.ProductID]
		if !ok {
			var err error
			menuItem, err = s.menuRepo.GetByID(item.ProductID)
			if err != nil {
				log.Print("Invalid product ID", "product_id", item.ProductID, "error", err)
				return models.Order{}, fmt.Errorf("product ID '%s' not found in menu", item.ProductID)
			}
			menuCache[item.ProductID] = menuItem
		}

		for _, ingredient := range menuItem.Ingredients {
			invItem, err := s.inventoryRepo.GetByID(ingredient.IngredientID)
			if err != nil {
				log.Printf("Inventory item not found", "ingredient_id", ingredient.IngredientID, "error", err)
				return models.Order{}, fmt.Errorf("ingredient '%s' not available", ingredient.IngredientID)
			}
			needed := ingredient.Quantity * (item.Quantity)
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

	order.Status = "open"

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Print("Failed to begin transaction", "error", err)
		return models.Order{}, errors.New("failed to start transaction")
	}

	rollback := func(err error) (models.Order, error) {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Print("Failed to rollback transaction", "error", rbErr)
		}
		return models.Order{}, err
	}

	for _, item := range order.Items {
		menuItem := menuCache[item.ProductID]
		for _, ingredient := range menuItem.Ingredients {
			invItem, err := s.inventoryRepo.GetByID(ingredient.IngredientID)
			if err != nil {
				return rollback(fmt.Errorf("ingredient '%s' not available", ingredient.IngredientID))
			}
			needed := ingredient.Quantity * item.Quantity
			invItem.Quantity -= needed
			updatedInv, err := s.inventoryRepo.UpdateTx(tx, invItem)
			if err != nil {
				log.Print("Failed to update inventory", "error", err)
				return rollback(fmt.Errorf("failed to update inventory: %v", err))
			}
			_ = updatedInv
		}
	}

	createdOrder, err := s.orderRepo.CreateTx(tx, order)
	if err != nil {
		log.Print("Failed to save order", "error", err)
		return rollback(errors.New("failed to save order"))
	}

	if err := tx.Commit(); err != nil {
		log.Print("Failed to commit transaction", "error", err)
		return models.Order{}, errors.New("failed to commit transaction")
	}

	return createdOrder, nil
}

func (s *orderService) GetOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *orderService) GetOrder(id int64) (models.Order, error) {
	if id == 0 {
		return models.Order{}, errors.New("id is required")
	}
	return s.orderRepo.GetByID(id)
}

func (s *orderService) UpdateOrder(id int64, order models.Order) (models.Order, error) {
	if id == 0 {
		return models.Order{}, errors.New("id is required")
	}

	existingOrder, err := s.orderRepo.GetByID(id)
	if err != nil {
		return models.Order{}, err
	}

	order.ID = existingOrder.ID
	order.CreatedAt = existingOrder.CreatedAt
	order.Status = existingOrder.Status

	return s.orderRepo.Update(id, order)
}

func (s *orderService) DeleteOrder(id int64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return s.orderRepo.Delete(id)
}

func (s *orderService) CloseOrder(id int64) (models.Order, error) {
	if id == 0 {
		return models.Order{}, errors.New("id is required")
	}

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

func (s *orderService) GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error) {
	return s.orderRepo.GetNumberOfOrderedItems(startDate, endDate)
}
