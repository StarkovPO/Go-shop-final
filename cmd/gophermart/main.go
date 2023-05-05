package main

import (
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	c, err := config.Init()
	if err != nil {
		logrus.Fatalf("init configuration: %s", err)
	}

	err = initApp(c)
	if err != nil {
		logrus.Fatalf("unsuccessful initilization app: %v", err)
	}

}
