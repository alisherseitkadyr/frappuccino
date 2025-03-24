package repository

import (
	"hot-coffee/models"
)

type InventoryRepository interface {
	GetAll() ([]models.InventoryItem, error)
	GetByID(id string) (models.InventoryItem, error)
	Update(item models.InventoryItem) (models.InventoryItem, error)
}

type inventoryRepository struct {
	store *FileStore
}

func NewInventoryRepository(dataDir string) InventoryRepository {
	return &inventoryRepository{
		store: NewFileStore(dataDir + "/inventory.json"),
	}
}

func (r *inventoryRepository) GetAll() ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	if err := r.store.Read(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *inventoryRepository) GetByID(id string) (models.InventoryItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}

	for _, item := range items {
		if item.IngredientID == id {
			return item, nil
		}
	}

	return models.InventoryItem{}, nil
}

func (r *inventoryRepository) Update(updatedItem models.InventoryItem) (models.InventoryItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}

	for i, item := range items {
		if item.IngredientID == updatedItem.IngredientID {
			items[i] = updatedItem
			if err := r.store.Write(items); err != nil {
				return models.InventoryItem{}, err
			}
			return updatedItem, nil
		}
	}

	return models.InventoryItem{}, nil
}
