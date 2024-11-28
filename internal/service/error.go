package service

import "errors"

var (
	ErrAccessSign     = errors.New("couldn't sign access token")
	ErrRefreshSign    = errors.New("couldn't sign refresh token")
	ErrRefreshInvalid = errors.New("refresh token is invalid")
	ErrRefreshExpired = errors.New("refresh token is expired")
	ErrWrongRefresh   = errors.New("provided refresh token doesn't match the stored one")
	ErrWrongIp        = errors.New("your ip adress doesn't match the stored one ")
	ErrMalformedToken = errors.New("your token is malformed")
)
