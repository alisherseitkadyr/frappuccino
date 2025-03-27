package service

import (
	"errors"
	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type InventoryService interface {
	CreateInventoryItem(name string, quantity float64, unit string) (models.InventoryItem, error)
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

func (s *inventoryService) CreateInventoryItem(name string, quantity float64, unit string) (models.InventoryItem, error) {
	if name == "" {
		return models.InventoryItem{}, errors.New("name is required")
	}
	if quantity < 0 {
		return models.InventoryItem{}, errors.New("quantity cannot be negative")
	}
	if unit == "" {
		return models.InventoryItem{}, errors.New("unit is required")
	}

	// Create new inventory item with generated ID
	item := models.NewInventoryItem(name, quantity, unit)

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
