package repository

import "errors"

var (
	ErrNotFound       = errors.New("data was not found")
	ErrNoEmail        = errors.New("email was not found")
	ErrNotImplemented = errors.New("not implemented")
)
