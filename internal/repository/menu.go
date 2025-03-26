package repository

import (
	"errors"
	"hot-coffee/models"
)

type MenuRepository interface {
	Create(item models.MenuItem) (models.MenuItem, error)
	GetAll() ([]models.MenuItem, error)
	GetByID(id string) (models.MenuItem, error)
	Update(id string, item models.MenuItem) (models.MenuItem, error)
	Delete(id string) error
}

type menuRepository struct {
	store *FileStore
}

func NewMenuRepository(dataDir string) MenuRepository {
	return &menuRepository{
		store: NewFileStore(dataDir + "/menu_items.json"),
	}
}

func (r *menuRepository) Create(item models.MenuItem) (models.MenuItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.MenuItem{}, err
	}

	// Check for duplicate ID
	for _, existing := range items {
		if existing.ID == item.ID {
			return models.MenuItem{}, ErrDuplicateID
		}
	}

	items = append(items, item)
	if err := r.store.Write(items); err != nil {
		return models.MenuItem{}, err
	}

	return item, nil
}

func (r *menuRepository) GetAll() ([]models.MenuItem, error) {
	var items []models.MenuItem
	if err := r.store.Read(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuRepository) GetByID(id string) (models.MenuItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.MenuItem{}, err
	}

	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}

	return models.MenuItem{}, ErrNotFound
}

func (r *menuRepository) Update(id string, updatedItem models.MenuItem) (models.MenuItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.MenuItem{}, err
	}

	for i, item := range items {
		if item.ID == id {
			// Preserve the ID from the path, not the body
			updatedItem.ID = id
			items[i] = updatedItem
			if err := r.store.Write(items); err != nil {
				return models.MenuItem{}, err
			}
			return updatedItem, nil
		}
	}

	return models.MenuItem{}, ErrNotFound
}

func (r *menuRepository) Delete(id string) error {
	items, err := r.GetAll()
	if err != nil {
		return err
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			return r.store.Write(items)
		}
	}

	return ErrNotFound
}

var (
	ErrDuplicateID = errors.New("duplicate ID")
	ErrNotFound    = errors.New("not found")
)
