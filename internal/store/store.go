package store

import (
	"context"
	"database/sql"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"time"
)

type Store struct {
	db Postgres
}

func NewStore(db Postgres) Store {
	return Store{db: db}
}

func (o *Store) CreateUserDB(ctx context.Context, user models.Users) error {

	timestamp := time.Now().Unix()

	stmt, err := o.db.db.PrepareContext(ctx, `
        INSERT INTO users (id, login, password_hash, created_at)
        VALUES ($1, $2, $3, to_timestamp($4))
    `)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, user.Id, user.Login, user.Password, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (o *Store) CheckLogin(ctx context.Context, login string) bool {

	var exist bool

	err := o.db.db.QueryRowContext(ctx, `
        SELECT EXISTS (SELECT 1 FROM users WHERE login = $1)
    `, login).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return exist
}

func (o *Store) GetUserPass(ctx context.Context, login string) (string, bool) {
	var hash string

	err := o.db.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE login = $1`, login).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
	}
	return hash, true
}
