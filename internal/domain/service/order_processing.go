package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"l0/internal/domain/models"
	"l0/internal/storage"
	"l0/internal/storage/cache"
	"l0/internal/transport/broker"
	"log/slog"
)

const (
	Subject                = "order"
	MessageRedeliveryCount = 1
)

type OrderProcessing struct {
	storage storage.Storage
	cache   cache.Cacher
	stream  *broker.Stan
}

func NewOrderProcessing(storage storage.Storage, cache cache.Cacher, stream *broker.Stan) *OrderProcessing {
	return &OrderProcessing{storage: storage, cache: cache, stream: stream}
}

func (o *OrderProcessing) Save() error {
	const fn = "service.OrderProcessing.Save"

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := o.stream.Subscribe(Subject, func(msg *stan.Msg) {
		if msg.Redelivered && msg.RedeliveryCount > MessageRedeliveryCount {
			err := msgAck(msg)
			if err != nil {
				slog.Error(fn, slog.String("failed to ack message error", err.Error()))
				return
			}
			return
		}

		var order models.Order
		data := msg.Data
		err := json.Unmarshal(data, &order)
		if err != nil {
			slog.Error(fn, slog.String("failed to unmarshal data error", err.Error()))
			return
		}

		err = validate.Struct(order)
		if err != nil {
			slog.Error(fn, slog.String("validate is fail", err.Error()))
			return
		}

		err = o.storage.SaveOrder(order)
		if err != nil {
			slog.Error(fn, slog.String("failed to save order error", err.Error()))
			return
		}

		o.cache.SaveOrderToCache(order)
		slog.Info("order is saved", slog.String(order.OrderUID, fmt.Sprint(order)))

		err = msgAck(msg)
		if err != nil {
			slog.Error(fn, slog.String("failed to ack message error: ", err.Error()))
			return
		}
		slog.Info("message is acked")
	})
	if err != nil {
		slog.Error(fn, slog.String("failed to subscribe error: ", err.Error()))
		return err
	}

	return nil
}

func (o *OrderProcessing) GetOrder(uid string) (models.Order, error) {
	return o.cache.GetOrderFromCache(uid)
}

func msgAck(msg *stan.Msg) error {
	err := msg.Ack()
	if err != nil {
		slog.Error("failed to ack message error: ", err)
		return err
	}
	slog.Info("message is acked")
	return nil
}
