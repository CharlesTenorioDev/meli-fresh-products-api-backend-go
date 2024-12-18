package errorss

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	StatusHttp int
	Message    string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.StatusHttp), e.Message)
}

func NewUnprocessableEntity(message string) error {
	return &CustomError{StatusHttp: http.StatusUnprocessableEntity, Message: message}
}

func NewBadRequestError(message string) error {
	return &CustomError{StatusHttp: http.StatusBadRequest, Message: message}
}

func NewConflictError(message string) error {
	return &CustomError{StatusHttp: http.StatusConflict, Message: message}
}

func NewNotFound(message string) error {
	return &CustomError{StatusHttp: http.StatusNotFound, Message: message}
}
