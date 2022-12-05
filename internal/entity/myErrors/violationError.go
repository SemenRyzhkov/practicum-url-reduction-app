package myErrors

import (
	"fmt"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type ViolationError struct {
	Err      error
	Response entity.URLResponse
}

func (ve *ViolationError) Error() string {
	return fmt.Sprintf("URL %s already exists", ve.Response.Result)
}

func (ve *ViolationError) Unwrap() error {
	return ve.Err
}

func NewViolationError(URL string, err error) error {
	response := entity.URLResponse{
		Result: URL,
	}
	return &ViolationError{
		Response: response,
		Err:      err,
	}
}
