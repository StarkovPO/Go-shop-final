package appErrors

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
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
