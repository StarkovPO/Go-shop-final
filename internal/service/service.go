package service

import (
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/StarkovPO/Go-shop-final/internal/store"
)

type Service struct {
	store  store.Store
	config config.Config
}

func NewService(s store.Store, c config.Config) Service {
	return Service{
		store:  s,
		config: c,
	}
}