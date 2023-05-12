package models

type Orders struct {
	UserID     string `json:"-" db:"-"`
	ID         int    `json:"number" db:"id"`
	Status     string `json:"status" db:"status"`
	Accrual    int    `json:"accrual" db:"accrual"`
	UploadedAt string `json:"uploaded_at" db:"uploaded_at"`
}
