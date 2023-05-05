package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"net/http"
)

var (
	appErr *appErrors.AppError
)

type ServiceInterface interface {
	CreateUser(ctx context.Context, req models.Users) (string, error)
	GenerateUserToken(ctx context.Context, req models.Users) (string, error)
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
			if errors.Is(err, appErrors.ErrLoginAlreadyExist) {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write(appErrors.ErrInvalidLoginOrPass.Marshal())
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

//func RegisterUser(s ServiceInterface) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//	}
//
//}
