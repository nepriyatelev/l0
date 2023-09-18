package models

import "time"

type Order struct {
	OrderUID        string    `json:"order_uid" validate:"required,eqcsfield=Payment.Transaction"`
	TrackNumber     string    `json:"track_number" validate:"required"`
	Entry           string    `json:"entry" validate:"required"`
	Delivery        Delivery  `json:"delivery" validate:"required"`
	Payment         Payment   `json:"payment" validate:"required"`
	Items           []Item    `json:"items" validate:"required"`
	Locale          string    `json:"locale" validate:"required"`
	InternalSig     string    `json:"internal_signature"`
	CustomerID      string    `json:"customer_id" validate:"required"`
	DeliveryService string    `json:"delivery_service" validate:"required"`
	ShardKey        string    `json:"shardkey" validate:"required"`
	SmID            int       `json:"sm_id" validate:"required"`
	DateCreated     time.Time `json:"date_created" validate:"required"`
	OofShard        string    `json:"oof_shard" validate:"required"`
}
