package service

import (
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	id guid.GUID
	ip string

	exp time.Time
	iat time.Time
}

func FromMap(claims jwt.MapClaims) tokenClaims {
	return tokenClaims{
		id:  claims["id"].(guid.GUID),
		ip:  claims["ip"].(string),
		exp: claims["exp"].(time.Time),
		iat: claims["iat"].(time.Time),
	}
}

func (tc tokenClaims) ToMap() jwt.MapClaims {
	return jwt.MapClaims{
		"ip":  tc.ip,
		"id":  tc.id,
		"exp": tc.exp.Unix(),
		"iat": tc.iat.Unix(),
	}
}

func generateToken(id guid.GUID, ip string, exp time.Duration) *jwt.Token {
	claims := tokenClaims{
		id:  id,
		ip:  ip,
		exp: time.Now().Add(exp),
		iat: time.Now(),
	}.ToMap()

	return jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		claims,
	)
}
