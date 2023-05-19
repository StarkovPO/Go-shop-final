package models

type Balance struct {
	PrimaryID string  `json:"-" db:"primary_id"`
	UserID    string  `json:"-" db:"user_id"`
	Current   float64 `json:"current" db:"current"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
}
