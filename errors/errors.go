package errors

import (
	"fmt"
)

// BusinessError 业务错误
type BusinessError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Error
func (err *BusinessError) Error() string {
	return fmt.Sprintf("%d - %s", err.Code, err.Msg)
}

// New ...
func New(code int, msg string) *BusinessError {
	return &BusinessError{
		Code: code,
		Msg:  msg,
	}
}
