package handler

import (
	"context"
	"encoding/json"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	appErr *appErrors.AppError
)

type ServiceInterface interface {
	CreateUser(ctx context.Context, req models.Users) (string, error)
	GenerateUserToken(ctx context.Context, req models.Users) (string, error)
	CreateUserOrder(ctx context.Context, req models.Orders) error
	GetUserOrders(ctx context.Context, UID string) ([]models.Orders, error)
	GetUserBalance(ctx context.Context, UID string) (models.Balance, error)
	CreateUserWithdraw(ctx context.Context, req models.Withdrawn) error
	GetUserWithdrawn(ctx context.Context, UID string) ([]models.Withdrawn, error)
}

func RegisterUser(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Users
		ctx := r.Context()

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		token, err := s.CreateUser(ctx, req)

		if err != nil {
			appErrors.HandleError(w, err)
			return
		}
		w.Header().Set("Authorization", token)
	}
}

func LoginUser(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Users
		ctx := r.Context()

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		token, err := s.GenerateUserToken(ctx, req)

		if err != nil {
			appErrors.HandleError(w, err)
			return
		}

		w.Header().Set("Authorization", token)
	}
}

func CreateOrder(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, appErrors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		body, _ := io.ReadAll(r.Body)
		UID := r.Header.Get("User-ID")
		ID := string(body)

		if ID == "" {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}
		logrus.Infof("Order ID in Handler: %v", ID)
		req := models.Orders{UserID: UID, ID: ID}

		err := s.CreateUserOrder(ctx, req)

		if err != nil {
			appErrors.HandleError(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func GetUserOrders(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, appErrors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		UID := r.Header.Get("User-ID")

		res, err := s.GetUserOrders(ctx, UID)

		if err != nil {
			appErrors.HandleError(w, err)
			return
		}

		w.Header().Set("Content-type", "application/json")

		b, err := json.Marshal(res)
		logrus.Infof("User orders: %v", string(b))
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			return
		}

	}
}

func GetUserBalance(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, appErrors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		UID := r.Header.Get("User-ID")

		res, err := s.GetUserBalance(ctx, UID)

		if err != nil {
			logrus.Errorf("ops something went wrong: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")

		b, err := json.Marshal(res)

		if err != nil {
			logrus.Errorf("ops something went wrong: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			return
		}
	}
}

func CreateUserWithdraw(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, appErrors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		var req models.Withdrawn

		ctx := r.Context()

		UID := r.Header.Get("User-ID")

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		req.UserID = UID
		err = s.CreateUserWithdraw(ctx, req)

		if err != nil {
			appErrors.HandleError(w, err)
			return
		}

	}
}

func GetUserWithdraw(s ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, appErrors.ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		UID := r.Header.Get("User-ID")

		res, err := s.GetUserWithdrawn(ctx, UID)

		if err != nil {
			appErrors.HandleError(w, err)
			return
		}

		w.Header().Set("Content-type", "application/json")

		b, err := json.Marshal(res)

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			return
		}
	}
}
