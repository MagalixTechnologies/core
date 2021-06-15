package error

import (
	"net/http"
)

type Error interface {
	Error() string
	Code() int
	Messages() []string
}

type CoreError struct {
	err      error
	httpCode int
	messages []string
}

func (e *CoreError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

// Code get error http code
func (e *CoreError) Code() int {
	return e.httpCode
}

// Message get error message
func (e *CoreError) Messages() []string {
	return e.messages
}

// ErrNotFound not found error
func ErrNotFound(messages ...string) error {
	return &CoreError{
		messages: messages,
		httpCode: http.StatusNotFound,
	}
}

// ErrInvalidInput invalid input error
func ErrInvalidInput(err error, messages ...string) error {
	return &CoreError{
		err:      err,
		messages: messages,
		httpCode: http.StatusBadRequest,
	}
}

// ErrConflict not found error
func ErrConflict(messages ...string) error {
	return &CoreError{
		messages: messages,
		httpCode: http.StatusConflict,
	}
}

// ErrOther generic internal server error
func ErrOther(err error, messages ...string) error {
	return &CoreError{
		err:      err,
		messages: messages,
		httpCode: http.StatusInternalServerError,
	}
}
