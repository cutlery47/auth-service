package service

import (
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
)

var (
	defaultAccessTTL  time.Duration = time.Hour
	defaultRefreshTTL time.Duration = time.Hour * 24
)

type tokenClaims struct {
	id guid.GUID
	ip string

	exp time.Time
	iat time.Time
	iss string
}

func NewClaims(id guid.GUID, ip string, exp time.Duration) *tokenClaims {
	return &tokenClaims{
		id:  id,
		ip:  ip,
		exp: time.Now().Add(exp),
		iat: time.Now(),
		iss: "auth-service",
	}
}
