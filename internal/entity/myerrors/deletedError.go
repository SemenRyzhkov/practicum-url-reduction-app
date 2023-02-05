package myerrors

import (
	"fmt"
)

// DeletedError ошибка после удаления
type DeletedError struct {
	Err error
	URL string
}

// Error ошибка в виде строки
func (de *DeletedError) Error() string {
	return fmt.Sprintf("URL %s is deleted", de.URL)
}

// Unwrap исходная ошибка
func (de *DeletedError) Unwrap() error {
	return de.Err
}

// NewDeletedError конструктор ошибки
func NewDeletedError(URL string, err error) error {
	return &DeletedError{
		URL: URL,
		Err: err,
	}
}
