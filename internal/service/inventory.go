package service

import (
	"errors"
	"frappuccino/internal/repository"
	"frappuccino/models"
)

type InventoryService interface {
	CreateInventoryItem(name string, quantity int, unit string) (models.InventoryItem, error)
	GetInventoryItems() ([]models.InventoryItem, error)
	GetInventoryItem(id int64) (models.InventoryItem, error)
	UpdateInventoryItem(item models.InventoryItem) (models.InventoryItem, error)
	DeleteInventoryItem(id int64) error

	GetLeftOvers(sortBy string, page, pageSize int) ([]models.InventoryItem, int, error)

}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) CreateInventoryItem(name string, quantity int, unit string) (models.InventoryItem, error) {
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

func (s *inventoryService) GetInventoryItem(id int64) (models.InventoryItem, error) {
	if id == 0 {
		return models.InventoryItem{}, errors.New("id is required")
	}
	return s.repo.GetByID(id)
}

func (s *inventoryService) DeleteInventoryItem(id int64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}


func (s *inventoryService) UpdateInventoryItem(item models.InventoryItem) (models.InventoryItem, error) {
	if item.IngredientID == 0 {
		return models.InventoryItem{}, errors.New("id is required")
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

	return s.repo.Update(item)
}


func (s *inventoryService) GetLeftOvers(sortBy string, page, pageSize int) ([]models.InventoryItem, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	items, totalCount, err := s.repo.GetLeftOvers(sortBy, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return items, totalCount, nil
}
