package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/StarkovPO/Go-shop-final/internal/apperrors"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/sirupsen/logrus"
	"time"
)

const duplicateUserBalance = "pq: duplicate key value violates unique constraint \"user_id_key\""

type Store struct {
	store Postgres
}

func NewStore(db Postgres) Store {
	return Store{store: db}
}

func (o *Store) CreateUserDB(ctx context.Context, user models.Users) error {

	timestamp := time.Now().Unix()

	stmt, err := o.store.db.PrepareContext(ctx, createUser)
	if err != nil {
		logrus.Errorf("unhandled error: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, user.ID, user.Login, user.Password, timestamp)
	if err != nil {
		logrus.Errorf("unhandled error: %v", err)
		return err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return nil
}

func (o *Store) CheckLogin(ctx context.Context, login string) bool {

	var exist bool

	stmt, err := o.store.db.PrepareContext(ctx, checkLogin)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
	}

	err = stmt.QueryRowContext(ctx, login).Scan(&exist)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("No rows returned")
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

	stmt, err := o.store.db.PrepareContext(ctx, getUserPass)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
	}

	err = stmt.QueryRowContext(ctx, login).Scan(&hash)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("No rows returned")
			return "", false
		}
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}
	return hash, true
}

func (o *Store) CreateUserOrderDB(ctx context.Context, order models.OrderFromService) error {

	var UID string

	timestamp := time.Now().Unix()

	stmt, err := o.store.db.PrepareContext(ctx, createOrder)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return err
	}

	stmt2, err := o.store.db.PrepareContext(ctx, getUserFromOrders)
	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, order.UserID, order.ID, order.Status, order.Accrual, timestamp)
	if err != nil {
		d := err.Error()
		if d == "pq: duplicate key value violates unique constraint \"orders_id_key\"" {
			err = stmt2.QueryRowContext(ctx, order.ID).Scan(&UID)
			if UID != order.UserID {
				logrus.Printf("Order belong to user %v, but got: %v", UID, order.UserID)
				return apperrors.ErrOrderAlreadyExist
			} else if UID == order.UserID {
				logrus.Printf("Order belong to user %v", order.UserID)
				return apperrors.ErrOrderAlreadyBelong
			} else {
				logrus.Errorf("ops unhandled error:%v", err)
				return err
			}
		}
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return nil
}

func (o *Store) GetUserOrders(ctx context.Context, UID string) ([]models.Orders, error) {

	var orders []models.Orders

	err := o.store.db.SelectContext(ctx, &orders, getOrders, UID)

	if err != nil {
		logrus.Errorf("unhandled error: %v", err)
		return nil, err
	}

	return orders, nil
}

func (o *Store) GetUserID(ctx context.Context, login string) (string, error) {
	var UID string

	stmt, err := o.store.db.PrepareContext(ctx, getUserID)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return "", err
	}

	err = stmt.QueryRowContext(ctx, login).Scan(&UID)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("impossible error: %v", err)
			return "", errors.New("login not found. Impossible")
		}
		logrus.Errorf("unhandled error: %v", err)
		return "", err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}
	return UID, nil
}

func (o *Store) IncreaseUserBalance(ctx context.Context, accrual float64, UID string) error {

	stmt, err := o.store.db.PrepareContext(ctx, createUserBalance)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, UID, accrual, 0)

	if err != nil {
		if err.Error() != duplicateUserBalance {
			logrus.Infof("handled error: %v", err)
			stmt, err = o.store.db.PrepareContext(ctx, increaseUserBalance)

			if err != nil {
				logrus.Errorf("error with stmt: %v", err)
				return err
			}

			_, err = stmt.ExecContext(ctx, UID, accrual)

			if err != nil {
				logrus.Errorf("ops unhandled error: %v", err)
				return err
			}
			return nil
		}
		logrus.Errorf("ops unhandled error: %v", err)
		return err
	}
	return nil

}

func (o *Store) DecreaseUserBalance(ctx context.Context, withdrawn float64, UID string) error {

	stmt, err := o.store.db.PrepareContext(ctx, decreaseUserBalance)

	if err != nil {
		logrus.Errorf("ops unhandled error: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, UID, withdrawn)

	if err != nil {
		logrus.Errorf("ops unhandled error: %v", err)
		return err
	}

	return err
}

func (o *Store) GetUserBalanceDB(ctx context.Context, UID string) (models.Balance, error) {
	var balance models.Balance

	err := o.store.db.GetContext(ctx, &balance, getUserBalance, UID)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			logrus.Infof("handled error: %v", err)
			return models.Balance{Current: 0, Withdrawn: 0}, nil
		}

		logrus.Errorf("unhandled error: %v", err)
		return models.Balance{}, err
	}
	return balance, nil

}

func (o *Store) IncreaseUserWithdrawn(ctx context.Context, withdrawn float64, UID string) error {

	stmt, err := o.store.db.PrepareContext(ctx, increaseUserWithdrawn)

	if err != nil {
		logrus.Errorf("ops unhandled error: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, UID, withdrawn)

	if err != nil {
		logrus.Errorf("ops unhandled error: %v", err)
		return err
	}

	return err
}

func (o *Store) CreateWithdraw(ctx context.Context, req models.Withdrawn) error {

	timestamp := time.Now().Unix()
	stmt, err := o.store.db.PrepareContext(ctx, createUserWithdrawn)

	if err != nil {
		logrus.Errorf("unexpected error: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, req.OrderID, req.Withdrawn, req.UserID, timestamp)

	if err != nil {
		logrus.Errorf("unexpected error: %v", err)
		return err
	}

	err = o.DecreaseUserBalance(ctx, req.Withdrawn, req.UserID)

	if err != nil {
		logrus.Errorf("unexpected error: %v", err)
		return err
	}

	err = o.IncreaseUserWithdrawn(ctx, req.Withdrawn, req.UserID)

	return err
}

func (o *Store) GetUserWithdrawnDB(ctx context.Context, UID string) ([]models.Withdrawn, error) {
	var withdrawn []models.Withdrawn

	err := o.store.db.SelectContext(ctx, &withdrawn, getUserWithdrawn, UID)

	if err != nil {
		logrus.Errorf("unhandled error: %v", err)
		return nil, err
	}

	return withdrawn, nil
}

func (o *Store) UpdateOrderStatus(ctx context.Context, status string, UID string) error {

	stmt, err := o.store.db.PrepareContext(ctx, updateOrderStatus)

	if err != nil {
		logrus.Errorf("ops unhandled error: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, status, UID)

	if err != nil {
		logrus.Errorf("ops unhandled error: %v", err)
		return err
	}

	return err
}
