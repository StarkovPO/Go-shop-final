package models

type Withdrawn struct {
	OrderID     string  `json:"order" db:"order_id"`
	Withdrawn   float64 `json:"sum" db:"withdrawn"`
	UserID      string  `json:"-"`
	ProcessedAt string  `json:"processed_at" db:"processed_at"`
}
