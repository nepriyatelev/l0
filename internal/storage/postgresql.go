package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"l0/internal/domain/models"
	"log/slog"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*PostgresStorage, error) {
	const fn = "storage.postgresql.NewStorage"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		slog.Error(fn, slog.String("failed to open db error", err.Error()))
		return nil, err
	}
	slog.Info("db is opened")

	err = db.Ping()
	if err != nil {
		slog.Error(fn, slog.String("failed to ping db error", err.Error()))
		return nil, err
	}
	slog.Info("db is pinged")

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Close() error {
	const fn = "storage.postgresql.Close"

	err := s.db.Close()
	if err != nil {
		slog.Error(fn, slog.String("failed to close db error", err.Error()))
		return err
	}
	slog.Info("db is closed")

	return nil
}

func (s *PostgresStorage) SaveOrder(order models.Order) error {
	const fn = "storage.postgresql.SaveOrderToCache"

	tx, err := s.db.Begin()
	if err != nil {
		slog.Error(fn, slog.String("failed to begin transaction error", err.Error()))
		return err
	}

	_, err = tx.Exec(`INSERT INTO orders
    									(order_uid,
    									 track_number,
    									 entry,
    									 locale,
    									 internal_signature,
    									 customer_id,
    									 delivery_service,
    									 shard_key,
    									 sm_id,
    									 date_created,
    									 oof_shard)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSig,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard)
	if err != nil {
		slog.Error(fn, slog.String("failed to insert order error", err.Error()))
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			slog.Error(fn, slog.String("failed to rollback transaction error", rollbackErr.Error()))
		}
		return err
	}
	slog.Info("order is inserted")

	_, err = tx.Exec(`INSERT INTO delivery
    									(order_uid,
    									 name,
    									 phone,
    									 zip,
    									 city,
    									 address,
    									 region,
    									 email)
    							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	if err != nil {
		slog.Error(fn, slog.String("failed to insert delivery error", err.Error()))
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			slog.Error(fn, slog.String("failed to rollback transaction error", rollbackErr.Error()))
		}
		return err
	}
	slog.Info("delivery is inserted")

	_, err = tx.Exec(`INSERT INTO payments (
--                       order_uid,
						transaction,
						request_id,
						currency,
						provider,
						amount,
						payment_dt,
						bank,
						delivery_cost,
						goods_total,
						custom_fee)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		//order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		slog.Error(fn, slog.String("failed to insert payment error", err.Error()))
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			slog.Error(fn, slog.String("failed to rollback transaction error", rollbackErr.Error()))
		}
		return err
	}
	slog.Info("payment is inserted")

	for _, item := range order.Items {
		_, err = tx.Exec(`INSERT INTO items (
--                       order_uid,
                   chrt_id,
                   track_number,
                   price,
                   rid,
                   name,
                   sale,
                   size,
                   total_price,
                   nm_id,
                   brand,
                   status)
    			   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			//order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status)
		if err != nil {
			slog.Error(fn, slog.String("failed to insert item error", err.Error()))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.Error(fn, slog.String("failed to rollback transaction error", rollbackErr.Error()))
			}
			return err
		}
		slog.Info("item is inserted")
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(fn, slog.String("failed to commit transaction error", err.Error()))
		return err
	}
	slog.Info("transaction is committed")
	return nil
}

func (s *PostgresStorage) GetAllOrders() ([]models.Order, error) {
	const fn = "storage.postgresql.GetAllOrders"
	orders, err := s.GetOrders()
	if err != nil {
		slog.Error(fn, slog.String("failed to get orders error", err.Error()))
		return nil, err
	}

	for i, order := range orders {
		delivery, deliveryErr := s.GetDelivery(order.OrderUID)
		if deliveryErr != nil {
			slog.Error(fn, slog.String("failed to get delivery error", deliveryErr.Error()))
			return nil, deliveryErr
		}

		payment, paymentErr := s.GetPayment(order.OrderUID)
		if paymentErr != nil {
			slog.Error(fn, slog.String("failed to get payment error", paymentErr.Error()))
			return nil, paymentErr
		}

		items, itemsErr := s.GetItems(order.TrackNumber)
		if itemsErr != nil {
			slog.Error(fn, slog.String("failed to get items error", itemsErr.Error()))
			return nil, itemsErr
		}

		orders[i].Delivery = delivery
		orders[i].Payment = payment
		orders[i].Items = items
	}
	slog.Info(fn, slog.String("orders", fmt.Sprint(orders)))
	return orders, nil
}

func (s *PostgresStorage) GetOrders() ([]models.Order, error) {
	const fn = "storage.postgresql.GetOrders"
	q := `SELECT * FROM orders`
	rows, err := s.db.Query(q)
	if err != nil {
		slog.Error(fn, slog.String("failed to get orders error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSig,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			slog.Error(fn, slog.String("failed to scan row error", err.Error()))
			return nil, err
		}
		slog.Info("row is scanned")
		slog.Info(fn, slog.String(order.OrderUID, fmt.Sprint(order))) // TODO: delete
		orders = append(orders, order)
	}
	slog.Info("orders are scanned")
	return orders, nil
}

func (s *PostgresStorage) GetDelivery(orderUid string) (models.Delivery, error) {
	const fn = "storage.postgresql.GetDelivery"
	q := `SELECT * FROM delivery WHERE delivery.order_uid = $1`
	rows, err := s.db.Query(q, orderUid)
	if err != nil {
		slog.Error(fn, slog.String("failed to get delivery error", err.Error()))
		return models.Delivery{}, err
	}
	defer rows.Close()
	var delivery models.Delivery
	for rows.Next() {
		err = rows.Scan(
			&delivery.OrderUID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			slog.Error(fn, slog.String("failed to scan row error", err.Error()))
			return models.Delivery{}, err
		}
		slog.Info("row is scanned")
		slog.Info(fn, slog.String(delivery.OrderUID, fmt.Sprint(delivery))) // TODO: delete
	}
	slog.Info("delivery is scanned")
	return delivery, nil
}

func (s *PostgresStorage) GetPayment(orderUid string) (models.Payment, error) {
	const fn = "storage.postgresql.GetPayment"
	q := `SELECT * FROM payments WHERE payments.transaction = $1`
	rows, err := s.db.Query(q, orderUid)
	if err != nil {
		slog.Error(fn, slog.String("failed to get payment error", err.Error()))
		return models.Payment{}, err
	}
	defer rows.Close()
	var payment models.Payment
	for rows.Next() {
		err = rows.Scan(
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDt,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			slog.Error(fn, slog.String("failed to scan row error", err.Error()))
			return models.Payment{}, err
		}
		slog.Info("row is scanned")
		slog.Info(fn, slog.String(payment.Transaction, fmt.Sprint(payment))) // TODO: delete
	}
	slog.Info("payment is scanned")
	return payment, nil
}

func (s *PostgresStorage) GetItems(trackNumber string) ([]models.Item, error) {
	const fn = "storage.postgresql.GetItems"
	q := `SELECT * FROM items WHERE items.track_number = $1`
	rows, err := s.db.Query(q, trackNumber)
	if err != nil {
		slog.Error(fn, slog.String("failed to get items error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		var item models.Item
		err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			slog.Error(fn, slog.String("failed to scan row error", err.Error()))
			return nil, err
		}
		slog.Info("row is scanned")
		slog.Info(fn, slog.String(item.TrackNumber, fmt.Sprint(item))) // TODO: delete
		items = append(items, item)
	}
	slog.Info("items are scanned")
	return items, nil
}
