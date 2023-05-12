package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"io"
	"net/http"
	"strconv"
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

		if errors.As(err, &appErr) {
			if errors.Is(err, appErrors.ErrLoginAlreadyExist) {
				w.WriteHeader(http.StatusConflict)
				_, err = w.Write(appErrors.ErrLoginAlreadyExist.Marshal())
				return
			} else if errors.Is(err, appErrors.ErrLoginAlreadyExist) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write(appErrors.ErrLoginAlreadyExist.Marshal())
				return
			}
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

		if errors.As(err, &appErr) {
			if errors.Is(err, appErrors.ErrInvalidLoginOrPass) {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write(appErrors.ErrInvalidLoginOrPass.Marshal())
				return
			}
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
		ID, err := strconv.Atoi(string(body)) // add validation for empty body

		if err != nil {
			http.Error(w, appErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		req := models.Orders{UserID: UID, ID: ID}

		err = s.CreateUserOrder(ctx, req)

		if errors.As(err, &appErr) {
			if errors.Is(err, appErrors.ErrInvalidLoginOrPass) {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write(appErrors.ErrInvalidLoginOrPass.Marshal())
				return
			} else if errors.Is(err, appErrors.ErrInvalidOrderNumber) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, err = w.Write(appErrors.ErrInvalidOrderNumber.Marshal())
				return
			} else if errors.Is(err, appErrors.ErrOrderAlreadyExist) {
				w.WriteHeader(http.StatusConflict)
				_, err = w.Write(appErrors.ErrOrderAlreadyExist.Marshal())
				return
			} else if errors.Is(err, appErrors.ErrOrderAlreadyBelong) {
				w.WriteHeader(http.StatusAccepted)
				_, err = w.Write(appErrors.ErrOrderAlreadyBelong.Marshal())
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		if errors.As(err, &appErr) {
			if errors.Is(err, appErrors.ErrInvalidLoginOrPass) {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write(appErrors.ErrInvalidLoginOrPass.Marshal())
				return
			} else if errors.Is(err, appErrors.ErrOrderNotFound) {
				w.WriteHeader(http.StatusNoContent)
				_, err = w.Write(appErrors.ErrOrderNotFound.Marshal())
				return
			}
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			http.Error(w, "Server error", http.StatusInternalServerError)
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
