package repository

import (
	"hot-coffee/models"
)

type InventoryRepository interface {
	Create(item models.InventoryItem) (models.InventoryItem, error)
	GetAll() ([]models.InventoryItem, error)
	GetByID(id string) (models.InventoryItem, error)
	Update(item models.InventoryItem) (models.InventoryItem, error)
	Delete(id string) error
}

type inventoryRepository struct {
	store *FileStore
}

func NewInventoryRepository(dataDir string) InventoryRepository {
	return &inventoryRepository{
		store: NewFileStore(dataDir + "/inventory.json"),
	}
}

func (r *inventoryRepository) Create(item models.InventoryItem) (models.InventoryItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}

	// Check for duplicate ID
	for _, existing := range items {
		if existing.IngredientID == item.IngredientID {
			return models.InventoryItem{}, ErrDuplicateID
		}
	}

	items = append(items, item)
	if err := r.store.Write(items); err != nil {
		return models.InventoryItem{}, err
	}

	return item, nil
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

	return models.InventoryItem{}, ErrNotFound
}

func (r *inventoryRepository) Update(item models.InventoryItem) (models.InventoryItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}

	for i, existing := range items {
		if existing.IngredientID == item.IngredientID {
			items[i] = item
			if err := r.store.Write(items); err != nil {
				return models.InventoryItem{}, err
			}
			return item, nil
		}
	}

	return models.InventoryItem{}, ErrNotFound
}

func (r *inventoryRepository) Delete(id string) error {
	items, err := r.GetAll()
	if err != nil {
		return err
	}

	for i, item := range items {
		if item.IngredientID == id {
			items = append(items[:i], items[i+1:]...)
			return r.store.Write(items)
		}
	}

	return ErrNotFound
}
