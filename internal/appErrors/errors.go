package appErrors

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

var (
	ErrBadRequest        = NewAppError(nil, "Bad request", "request in not supported format")
	ErrLoginAlreadyExist = NewAppError(nil, "Login already exist", "User with this login already exist in DB")
	ErrCreateUser        = NewAppError(nil, "Something went wrong", "Executing query to create user ")
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
