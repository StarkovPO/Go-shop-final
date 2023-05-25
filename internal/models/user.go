package models

type Users struct {
	ID       string `json:"-" db:"id"`
	Login    string `json:"login" binding:"require" db:"login"`
	Password string `json:"password" binding:"require" db:"password_hash"`
}
