package service

import (
	"errors"
	"frappuccino/internal/repository"
	"frappuccino/models"
)

type MenuService interface {
	CreateMenuItem(item models.MenuItem) (models.MenuItem, error)
	GetMenuItems() ([]models.MenuItem, error)
	GetMenuItem(id int64) (models.MenuItem, error)
	UpdateMenuItem(id int64, item models.MenuItem) (models.MenuItem, error)
	DeleteMenuItem(id int64) error
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) CreateMenuItem(item models.MenuItem) (models.MenuItem, error) {
	if item.Name == "" {
		return models.MenuItem{}, errors.New("name is required")
	}
	if item.Price <= 0 {
		return models.MenuItem{}, errors.New("price must be positive")
	}
	return s.repo.Create(item)
}

func (s *menuService) GetMenuItems() ([]models.MenuItem, error) {
	return s.repo.GetAll()
}

func (s *menuService) GetMenuItem(id int64) (models.MenuItem, error) {
	if id == 0 {
		return models.MenuItem{}, errors.New("id is required")
	}
	return s.repo.GetByID(id)
}

func (s *menuService) DeleteMenuItem(id int64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}

func (s *menuService) UpdateMenuItem(id int64, item models.MenuItem) (models.MenuItem, error) {
	if id != item.ID {
		return models.MenuItem{}, errors.New("ID in path doesn't match ID in body")
	}
	if item.Name == "" {
		return models.MenuItem{}, errors.New("name is required")
	}
	if item.Price <= 0 {
		return models.MenuItem{}, errors.New("price must be positive")
	}
	return s.repo.Update(id, item)
}
