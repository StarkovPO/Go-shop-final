package store

import (
	"fmt"
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const (
	DBMaxOpenConnection     = 25
	DBMaxIdleConnection     = 25
	DBMaxConnectionLifeTime = 10 * time.Minute
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func MustPostgresConnection(c config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", c.DatabaseURIValue)
	if err != nil {
		panic(err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		defer db.Close()
		return nil
	}

	db.SetMaxOpenConns(DBMaxOpenConnection)
	db.SetMaxIdleConns(DBMaxIdleConnection)
	db.SetConnMaxLifetime(DBMaxConnectionLifeTime)

	if err = MakeDB(*db); err != nil {
		panic(err)
	}

	return db
}

func MakeDB(db sqlx.DB) error {

	if _, err := db.Exec(createUserTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}
	if _, err := db.Exec(createOrderTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}
	if _, err := db.Exec(createBalanceTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}
	if _, err := db.Exec(createWithdrawTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := db.Exec(createBalanceForeignKey); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}
	if _, err := db.Exec(createOrderForeignKey); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}
	if _, err := db.Exec(createWithdrawForeignKey); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := db.Exec(createLoginIndex); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := db.Exec(createOrderIDIndex); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := db.Exec(createUserIDIndexBalance); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	return nil
}
