package main

import (
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"log"
)

func main() {
	c, err := config.Init()
	if err != nil {
		log.Fatalf("init configuration: %s", err)
	}

	err = initApp(c)

}
