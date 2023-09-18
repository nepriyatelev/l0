package cache

import (
	"errors"
	"fmt"
	"l0/internal/domain/models"
	"l0/internal/storage"
	"log/slog"
	"strconv"
	"sync"
)

type MemoryCash struct {
	orders map[string]models.Order
	mut    sync.RWMutex
}

func NewMemoryCash(db storage.GetAllOrderer) (*MemoryCash, error) {
	const fn = "storage.cache.NewMemoryCash"
	orders, err := db.GetAllOrders()
	if err != nil {
		slog.Error(fn, slog.String("failed to get all orders error: ", err.Error()))
		return nil, err
	}
	slog.Info("got all orders from db", slog.String("orders", fmt.Sprint(orders)))
	mapOrders := make(map[string]models.Order)
	for _, order := range orders {
		mapOrders[order.OrderUID] = order
		slog.Info("order is added to cache", slog.String(order.OrderUID, fmt.Sprint(order)))
	}
	slog.Info("cache is ready and filled", slog.String("orders", strconv.Itoa(len(mapOrders))))

	return &MemoryCash{orders: mapOrders}, nil
}

func (m *MemoryCash) SaveOrderToCache(order models.Order) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.orders[order.OrderUID] = order
	slog.Info("order successfully saved to cache", slog.String(order.OrderUID, fmt.Sprint(order)))
}

func (m *MemoryCash) GetOrderFromCache(uid string) (models.Order, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	const fn = "storage.cache.GetOrderFromCache"
	var order models.Order
	order, ok := m.orders[uid]
	if !ok {
		slog.Error(fn, slog.String("failed to find order in cache error: ", ErrOrderNotFound))
		return models.Order{}, errors.New(ErrOrderNotFound)
	}
	slog.Info("order successfully found in cache", slog.String("order.OrderUID", fmt.Sprint(order)))
	return order, nil
}
