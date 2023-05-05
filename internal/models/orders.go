package models

type Orders struct {
	UserID    string `json:"-"`
	ID        int    `json:"number"`
	Status    string `json:"status"`
	Accrual   int    `json:"accrual"`
	UpdatedAt string `json:"updated_at"`
}
