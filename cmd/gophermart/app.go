package main

import (
	"context"
	"fmt"
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/StarkovPO/Go-shop-final/internal/handler"
	"github.com/StarkovPO/Go-shop-final/internal/middleware"
	"github.com/StarkovPO/Go-shop-final/internal/service"
	"github.com/StarkovPO/Go-shop-final/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func initApp(c config.Config) error {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	db := store.NewPostgres(store.MustPostgresConnection(c))
	storeApp := store.NewStore(*db)
	serviceApp := service.NewService(storeApp, c)
	router := setupAPI(serviceApp)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	listener, err := net.Listen("tcp", c.RunAddressValue)
	if err != nil {
		logrus.Fatalf("failed to listen on address %v: %s", c.RunAddressValue, err.Error())
	}

	server := &http.Server{
		Handler:           router,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	go func() {
		if err = server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server error: %v", err)
		}
	}()
	logrus.Infof("Server started and listening on address and port %s", c.RunAddressValue)

	sig := <-cancelChan

	logrus.Infof("Caught signal %v", sig)
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %v", err)
	}

	logrus.Info("Server shutdown successfully")
	return nil
}

func setupAPI(s service.Service) *mux.Router {

	router := mux.NewRouter()
	router.Use(middleware.CheckToken)
	router.HandleFunc("/api/user/register", handler.RegisterUser(&s)).Methods(http.MethodPost)
	router.HandleFunc("/api/user/login", handler.LoginUser(&s)).Methods(http.MethodPost)

	return router

}
