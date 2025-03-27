package service

import (
	"errors"
	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type MenuService interface {
	CreateMenuItem(item models.MenuItem) (models.MenuItem, error)
	GetMenuItems() ([]models.MenuItem, error)
	GetMenuItem(id string) (models.MenuItem, error)
	UpdateMenuItem(id string, item models.MenuItem) (models.MenuItem, error)
	DeleteMenuItem(id string) error
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
	item.ID = generateOrderID()
	return s.repo.Create(item)
}

func (s *menuService) GetMenuItems() ([]models.MenuItem, error) {
	return s.repo.GetAll()
}

func (s *menuService) GetMenuItem(id string) (models.MenuItem, error) {
	return s.repo.GetByID(id)
}

func (s *menuService) UpdateMenuItem(id string, item models.MenuItem) (models.MenuItem, error) {
	if id != item.ID {
		return models.MenuItem{}, errors.New("ID in path doesn't match ID in body")
	}
	return s.repo.Update(id, item)
}

func (s *menuService) DeleteMenuItem(id string) error {
	return s.repo.Delete(id)
}
