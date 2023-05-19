package apperrors

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrBadRequest         = NewAppError(nil, "Bad request", "request in not supported format")
	ErrLoginAlreadyExist  = NewAppError(nil, "Login already exist", "User with this login already exist in DB")
	ErrCreateUser         = NewAppError(nil, "Something went wrong", "Executing query to create user ")
	ErrInvalidAuthHeader  = NewAppError(nil, "Invalid auth header", "Token do not contain Bearer or body")
	ErrInvalidLoginOrPass = NewAppError(nil, "Invalid login or password", "User tried to login via incorrect login or pass")
	ErrInvalidOrderNumber = NewAppError(nil, "Invalid or empty order number", "Order number is not valid. Checked by Luhn")
	ErrOrderAlreadyExist  = NewAppError(nil, "Order already belong to another user", "User tried to connect existed user's order")
	ErrOrderAlreadyBelong = NewAppError(nil, "Order already belong to current user", "User tried to connect again his order")
	ErrOrderNotFound      = NewAppError(nil, "User's order not found", "Nothing returned from DB")
	ErrNotEnoughPoints    = NewAppError(nil, "Not enough to write off points", "User send tried to write off more that have")
	ErrWithdrawnNotFound  = NewAppError(nil, "User withdrawn not found", "DB returned 0 row")
)

type AppError struct {
	Err    error  `json:"-"`
	Msg    string `json:"msg"`
	DevMsg string `json:"dev_msg"`
}

func NewAppError(err error, msg, devMsg string) *AppError {
	logrus.Errorf("Error: %v", devMsg)
	return &AppError{
		err,
		msg,
		devMsg,
	}
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBadRequest):
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrLoginAlreadyExist):
		http.Error(w, ErrLoginAlreadyExist.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrCreateUser):
		http.Error(w, ErrCreateUser.Error(), http.StatusInternalServerError)
	case errors.Is(err, ErrInvalidAuthHeader):
		http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrInvalidLoginOrPass):
		http.Error(w, ErrInvalidLoginOrPass.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrInvalidOrderNumber):
		http.Error(w, ErrInvalidOrderNumber.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, ErrOrderAlreadyExist):
		http.Error(w, ErrOrderAlreadyExist.Error(), http.StatusConflict)
	case errors.Is(err, ErrOrderAlreadyBelong):
		http.Error(w, ErrOrderAlreadyBelong.Error(), http.StatusOK)
	case errors.Is(err, ErrOrderNotFound):
		http.Error(w, ErrOrderNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, ErrNotEnoughPoints):
		http.Error(w, ErrNotEnoughPoints.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrWithdrawnNotFound):
		http.Error(w, ErrWithdrawnNotFound.Error(), http.StatusNotFound)
	default:
		logrus.Errorf("Unhandled error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
