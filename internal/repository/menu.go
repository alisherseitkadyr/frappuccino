package repository

import (
	"hot-coffee/models"
)

type MenuRepository interface {
	GetAll() ([]models.MenuItem, error)
	GetByID(id string) (models.MenuItem, error)
}

type menuRepository struct {
	store *FileStore
}

func NewMenuRepository(dataDir string) MenuRepository {
	return &menuRepository{
		store: NewFileStore(dataDir + "/menu_items.json"),
	}
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

	return models.MenuItem{}, nil
}
