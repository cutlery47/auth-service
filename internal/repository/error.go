package repository

import "errors"

var (
	ErrUserNotFound   = errors.New("user was not found")
	ErrNoEmail        = errors.New("email was not found")
	ErrNotImplemented = errors.New("not implemented")
)
