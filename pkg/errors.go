package pkg

import (
	"errors"
	"fmt"
)

var ErrLackSessionID = errors.New("lack sessionID")

var ErrLackResponseChan = errors.New("lack response chan")

type ResponseError struct {
	Code    int
	Message string
	Data    interface{}
}

func NewResponseError(code int, message string, data interface{}) *ResponseError {
	return &ResponseError{Code: code, Message: message, Data: data}
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("code=%d message=%s data=%+v", r.Code, r.Message, r.Data)
}
