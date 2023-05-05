package main

import (
	"context"
	"fmt"
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/StarkovPO/Go-shop-final/internal/handler"
	"github.com/StarkovPO/Go-shop-final/internal/service"
	"github.com/StarkovPO/Go-shop-final/internal/store"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func initApp(c config.Config) error {

	db := store.NewPostgres(store.MustPostgresConnection(c))
	s := service.NewService(db, c)
	router := setupAPI(s)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	listener, err := net.Listen("tcp", c.RunAddressValue)
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %v", c.RunAddressValue, err)
	}

	server := &http.Server{
		Handler:           router,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	go func() {
		if err = server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server started and listening on address and port %s", c.RunAddressValue)

	sig := <-cancelChan

	log.Printf("Caught signal %v", sig)
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %v", err)
	}

	log.Println("Server shutdown successfully")
	return nil
}

func setupAPI(s service.ServiceInteface) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", handler.Test(s)).Methods(http.MethodGet)

	return router

}
