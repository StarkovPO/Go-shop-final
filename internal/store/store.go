package store

import (
	"context"
	"database/sql"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/sirupsen/logrus"
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

	stmt, err := o.db.db.PrepareContext(ctx, createUser)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, user.Id, user.Login, user.Password, timestamp)
	if err != nil {
		return err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return nil
}

func (o *Store) CheckLogin(ctx context.Context, login string) bool {

	var exist bool

	stmt, err := o.db.db.PrepareContext(ctx, checkLogin)

	err = stmt.QueryRowContext(ctx, login).Scan(&exist)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return exist
}

func (o *Store) GetUserPass(ctx context.Context, login string) (string, bool) {
	var hash string

	stmt, err := o.db.db.PrepareContext(ctx, getUserPass)

	err = stmt.QueryRowContext(ctx, login).Scan(&hash)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}
	return hash, true
}

func (o *Store) CreateUserOrderDB(ctx context.Context, order models.Orders) error {

	var UID string

	timestamp := time.Now().Unix()

	stmt, err := o.db.db.PrepareContext(ctx, createOrder)
	stmt2, err := o.db.db.PrepareContext(ctx, getUserFromOrders)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, order.UserID, order.ID, order.Status, order.Accrual, timestamp)
	if err != nil {
		d := err.Error()
		if d == "pq: duplicate key value violates unique constraint \"orders_id_key\"" {
			err = stmt2.QueryRowContext(ctx, order.ID).Scan(&UID)
			if UID != order.UserID {
				return appErrors.ErrOrderAlreadyExist
			} else if UID == order.UserID {
				return appErrors.ErrOrderAlreadyBelong
			} else {
				return err
			}
		}
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return nil
}
