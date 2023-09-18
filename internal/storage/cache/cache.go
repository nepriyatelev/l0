package cache

import (
	"l0/internal/domain/models"
)

const (
	ErrOrderNotFound = "order not found"
)

type Cacher interface {
	SaveOrderedCacher
	GetOrderedCacher
}

type SaveOrderedCacher interface {
	SaveOrderToCache(orders models.Order)
}

type GetOrderedCacher interface {
	GetOrderFromCache(uid string) (models.Order, error)
}
