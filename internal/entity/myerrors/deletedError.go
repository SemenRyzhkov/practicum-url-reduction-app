package myerrors

import (
	"fmt"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type DeletedError struct {
	Err      error
	Response entity.URLDTO
}

func (ve *DeletedError) Error() string {
	return fmt.Sprintf("URL %s is deleted", ve.Response.OriginalURL)
}

func (ve *DeletedError) Unwrap() error {
	return ve.Err
}

func NewDeletedError(URL entity.URLDTO, err error) error {
	return &DeletedError{
		Response: URL,
		Err:      err,
	}
}
