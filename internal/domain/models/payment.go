package models

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
	PaymentDt    int64   `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"required"`
	GoodsTotal   float64 `json:"goods_total" validate:"required"`
	CustomFee    float64 `json:"custom_fee"`
}
