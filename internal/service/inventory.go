package service

import (
	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type InventoryService interface {
	GetInventoryItems() ([]models.InventoryItem, error)
	GetInventoryItem(id string) (models.InventoryItem, error)
	UpdateInventoryItem(item models.InventoryItem) (models.InventoryItem, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
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
