package myerrors

import (
	"fmt"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

// ViolationError ошибка, вызванная сохранением повторяющейся сущности
type ViolationError struct {
	Err      error
	Response entity.URLResponse
}

// Error ошибка в виде строки
func (ve *ViolationError) Error() string {
	return fmt.Sprintf("URL %s already exists", ve.Response.Result)
}

// Unwrap исходная ошибка
func (ve *ViolationError) Unwrap() error {
	return ve.Err
}

// NewViolationError конструктор ошибки
func NewViolationError(URL string, err error) error {
	response := entity.URLResponse{
		Result: URL,
	}
	return &ViolationError{
		Response: response,
		Err:      err,
	}
}
