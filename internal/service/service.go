package service

import (
	"fmt"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/config"
	"github.com/cutlery47/auth-service/internal/models"
	"github.com/cutlery47/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Auth(id guid.GUID, ip string) (access, refresh string, err error)
	Refresh(id guid.GUID, ip, refresh string) (newAccess, newRefresh string, err error)
}

type AuthService struct {
	repo repository.Repository
	conf config.Service
}

func (as *AuthService) Auth(id guid.GUID, ip string) (access, refresh string, err error) {
	accessClaims := NewClaims(id, ip, as.conf.AccessTTL)
	refreshClaims := NewClaims(id, ip, as.conf.RefreshTTL)

	accessToken := as.generateToken(*accessClaims)
	refreshToken := as.generateToken(*refreshClaims)

	signedAccess, err := accessToken.SignedString(as.conf.Secret)
	if err != nil {
		return "", "", fmt.Errorf("accessToken.SignedString: %v", err)
	}

	signedRefresh, err := refreshToken.SignedString(as.conf.Secret)
	if err != nil {
		return "", "", fmt.Errorf("refreshToken.SignedString: %v", err)
	}

	salt := uuid.New()
	saltedRefresh := fmt.Sprintf("%v:%v", salt, signedRefresh)

	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(saltedRefresh), as.conf.Cost)
	if err != nil {
		return "", "", fmt.Errorf("bcrypt.GenerateFromPassword: %v", err)
	}

	inRefresh := models.InRefresh{
		UserId: id,
		Hash:   hashedRefresh,
		Salt:   salt,
		Cost:   as.conf.Cost,
	}

	if err := as.repo.Create(inRefresh); err != nil {
		return "", "", err
	}

	return signedAccess, signedRefresh, nil
}

func (as *AuthService) Refresh(id guid.GUID, ip, refresh string) (newAccess, newRefresh string, err error) {
	return "", "", nil
}

func (as *AuthService) generateToken(claims tokenClaims) *jwt.Token {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"ip":  claims.ip,
			"id":  claims.id,
			"exp": claims.exp,
			"iat": claims.iat,
			"iss": claims.iss,
		},
	)
}
