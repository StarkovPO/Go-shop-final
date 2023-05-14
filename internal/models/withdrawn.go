package models

type Withdrawn struct {
	OrderID     string  `json:"order"`
	Withdrawn   float64 `json:"sum"`
	UserID      string  `json:"-"`
	ProcessedAt string  `json:"processed_at"`
}
