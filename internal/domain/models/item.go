package models

type Item struct {
	ChrtID      int     `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Rid         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        float64 `json:"sale" validate:"required"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float64 `json:"total_price" validate:"required"`
	NmID        int     `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      int     `json:"status" validate:"required"`
}
