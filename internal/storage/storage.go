package storage

import (
	"l0/internal/domain/models"
)

type Storage interface {
	SaveOrderer
	GetAllOrderer
}

type SaveOrderer interface {
	SaveOrder(orders models.Order) error
}

type GetAllOrderer interface {
	GetAllOrders() ([]models.Order, error)
}
