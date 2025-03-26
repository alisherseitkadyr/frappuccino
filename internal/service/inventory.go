package service

import (
	"errors"
	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type InventoryService interface {
	CreateInventoryItem(item models.InventoryItem) (models.InventoryItem, error)
	GetInventoryItems() ([]models.InventoryItem, error)
	GetInventoryItem(id string) (models.InventoryItem, error)
	UpdateInventoryItem(item models.InventoryItem) (models.InventoryItem, error)
	DeleteInventoryItem(id string) error
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) CreateInventoryItem(item models.InventoryItem) (models.InventoryItem, error) {
	if item.IngredientID == "" {
		return models.InventoryItem{}, errors.New("ingredient_id is required")
	}
	if item.Name == "" {
		return models.InventoryItem{}, errors.New("name is required")
	}
	if item.Quantity < 0 {
		return models.InventoryItem{}, errors.New("quantity cannot be negative")
	}
	if item.Unit == "" {
		return models.InventoryItem{}, errors.New("unit is required")
	}

	return s.repo.Create(item)
}

func (s *inventoryService) GetInventoryItems() ([]models.InventoryItem, error) {
	return s.repo.GetAll()
}

func (s *inventoryService) GetInventoryItem(id string) (models.InventoryItem, error) {
	return s.repo.GetByID(id)
}

func (s *inventoryService) UpdateInventoryItem(item models.InventoryItem) (models.InventoryItem, error) {
	return s.repo.Update(item)
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	return s.repo.Delete(id)
}
