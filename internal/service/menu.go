package service

import (
	"hot-coffee/internal/repository"
	"hot-coffee/models"
)

type MenuService interface {
	GetMenuItems() ([]models.MenuItem, error)
	GetMenuItem(id string) (models.MenuItem, error)
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) GetMenuItems() ([]models.MenuItem, error) {
	return s.repo.GetAll()
}

func (s *menuService) GetMenuItem(id string) (models.MenuItem, error) {
	return s.repo.GetByID(id)
}
