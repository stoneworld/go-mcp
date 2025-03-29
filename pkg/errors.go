package pkg

import (
	"errors"
	"fmt"
)

var ErrServerNotSupport = errors.New("this feature serve not support")

var ErrLackResponseChan = errors.New("lack response chan")

type LackSessionError struct {
	SessionID string
}

func NewLackSessionError(sessionID string) *LackSessionError {
	return &LackSessionError{SessionID: sessionID}
}

func (e *LackSessionError) Error() string {
	return fmt.Sprintf("lack session, sessionID=%+v", e.SessionID)
}

type ResponseError struct {
	Code    int
	Message string
	Data    interface{}
}

func NewResponseError(code int, message string, data interface{}) *ResponseError {
	return &ResponseError{Code: code, Message: message, Data: data}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("code=%d message=%s data=%+v", e.Code, e.Message, e.Data)
}
