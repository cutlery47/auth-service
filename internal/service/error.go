package service

import "errors"

var (
	ErrAccessSign  = errors.New("couldn't sign access token")
	ErrRefreshSign = errors.New("couldn't sign refresh token")
)
