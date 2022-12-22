package myerrors

import (
	"fmt"
)

type DeletedError struct {
	Err error
	URL string
}

func (de *DeletedError) Error() string {
	return fmt.Sprintf("URL %s is deleted", de.URL)
}

func (de *DeletedError) Unwrap() error {
	return de.Err
}

func NewDeletedError(URL string, err error) error {
	return &DeletedError{
		URL: URL,
		Err: err,
	}
}
