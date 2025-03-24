package repository

import (
	"hot-coffee/models"
)

type OrderRepository interface {
	Create(order models.Order) (models.Order, error)
	GetAll() ([]models.Order, error)
	GetByID(id string) (models.Order, error)
	Update(id string, order models.Order) (models.Order, error)
	Delete(id string) error
}

type orderRepository struct {
	store *FileStore
}

func NewOrderRepository(dataDir string) OrderRepository {
	return &orderRepository{
		store: NewFileStore(dataDir + "/orders.json"),
	}
}

func (r *orderRepository) Create(order models.Order) (models.Order, error) {
	orders, err := r.GetAll()
	if err != nil {
		return models.Order{}, err
	}

	orders = append(orders, order)
	if err := r.store.Write(orders); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) GetAll() ([]models.Order, error) {
	var orders []models.Order
	if err := r.store.Read(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetByID(id string) (models.Order, error) {
	orders, err := r.GetAll()
	if err != nil {
		return models.Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}

	return models.Order{}, nil
}

func (r *orderRepository) Update(id string, updatedOrder models.Order) (models.Order, error) {
	orders, err := r.GetAll()
	if err != nil {
		return models.Order{}, err
	}

	for i, order := range orders {
		if order.ID == id {
			orders[i] = updatedOrder
			if err := r.store.Write(orders); err != nil {
				return models.Order{}, err
			}
			return updatedOrder, nil
		}
	}

	return models.Order{}, nil
}

func (r *orderRepository) Delete(id string) error {
	orders, err := r.GetAll()
	if err != nil {
		return err
	}

	for i, order := range orders {
		if order.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			return r.store.Write(orders)
		}
	}

	return nil
}
