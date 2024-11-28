package service

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/config"
	"github.com/cutlery47/auth-service/internal/models"
	"github.com/cutlery47/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, id guid.GUID, ip string) (access, refresh string, err error)
	Refresh(ctx context.Context, id guid.GUID, ip, refresh string) (newAccess, newRefresh string, err error)
}

type AuthService struct {
	repo repository.Repository
	auth smtp.Auth

	conf config.Service
}

func NewAuthService(repo repository.Repository, conf config.Service) *AuthService {
	// при инициализации сервиса производим smtp-авторизацию
	auth := smtp.PlainAuth(
		"",
		conf.SMTP.Username,
		conf.SMTP.Password,
		conf.SMTP.Hostname,
	)

	return &AuthService{
		repo: repo,
		auth: auth,
		conf: conf,
	}
}

func (as *AuthService) Create(ctx context.Context, id guid.GUID, ip string) (access, refresh string, err error) {
	// создаем access и refresh JWT-токены
	accessToken := generateToken(id, ip, as.conf.AccessTTL)
	refreshToken := generateToken(id, ip, as.conf.RefreshTTL)

	// подписываем созданные токены
	signedAccess, err := accessToken.SignedString([]byte(as.conf.Secret))
	if err != nil {
		return "", "", fmt.Errorf("accessToken.SignedString: %v", err)
	}

	signedRefresh, err := refreshToken.SignedString([]byte(as.conf.Secret))
	if err != nil {
		return "", "", fmt.Errorf("refreshToken.SignedString: %v", err)
	}

	// добавляем соль к refresh-токену
	salt := uuid.New()
	saltedRefresh := fmt.Sprintf("%v:%v", salt, signedRefresh)

	// так как bcrypt принимает строки длины до 72 символов,
	// хэшиурем засоленный refresh-токен
	hasher := sha512.New()
	if _, err := hasher.Write([]byte(saltedRefresh)); err != nil {
		return "", "", fmt.Errorf("hasher.Write: %v", err)
	}
	truncatedRefresh := hasher.Sum(nil)

	// хэшируем refresh-токен bcrypt-ом
	hashedRefresh, err := bcrypt.GenerateFromPassword(truncatedRefresh, as.conf.Cost)
	if err != nil {
		return "", "", fmt.Errorf("bcrypt.GenerateFromPassword: %v", err)
	}

	inRefresh := models.InRefresh{
		UserId: id,
		Hash:   hashedRefresh,
		Salt:   salt,
		Cost:   as.conf.Cost,
	}

	// заносим информацию о новом refresh-токене в хранилище
	if err := as.repo.Create(ctx, inRefresh); err != nil {
		return "", "", err
	}

	return signedAccess, signedRefresh, nil
}

func (as *AuthService) Refresh(ctx context.Context, id guid.GUID, ip, refresh string) (newAccess, newRefresh string, err error) {
	refreshClaims := jwt.MapClaims{}

	// Парсинг и валидация refresh-токена
	_, err = jwt.ParseWithClaims(refresh, refreshClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(as.conf.Secret), nil
	})
	if err != nil {
		switch {
		// проверка на время жизни токена
		case errors.Is(err, jwt.ErrTokenExpired):
			return "", "", ErrRefreshExpired
		// проверка на структуру токена
		case errors.Is(err, jwt.ErrTokenMalformed):
			return "", "", ErrMalformedToken
		default:
			return "", "", fmt.Errorf("jwt.ParseWithClaims: %v", err)
		}

	}

	refreshId, _ := refreshClaims["id"].(string)
	refreshIp, _ := refreshClaims["ip"].(string)

	// Сверяем ip-адрес внутри токена с адресом, с которого поступил запрос
	// В случае несовпадения - посылаем письмо на почту пользователя, id которого указан в токене
	// В принципе, можно было отслылать письма асинхронно, но в данном случае решил не заморачиваться
	if ip != refreshIp {
		if err := as.sendWarning(ctx, refreshId, ip); err != nil {
			return "", "", err
		}
		return "", "", ErrWrongIp
	}

	storedRefresh, err := as.repo.Get(ctx, id)
	if err != nil {
		return "", "", err
	}

	// Проверяем сходится ли токен, полученный от пользователя, с токеном, лежащим в БД
	// Для этого нам нужно:
	// 1) К полученному токену добавить ту же соль, которая использовалась изначально
	// 2) Захэшировать полученный токен при помощи SHA512
	// 3) Сравнить bcrypt-хэши

	// Добавляем соль
	saltedRefresh := fmt.Sprintf("%v:%v", storedRefresh.Salt, refresh)

	// SHA512-хэш
	hasher := sha512.New()
	if _, err := hasher.Write([]byte(saltedRefresh)); err != nil {
		return "", "", fmt.Errorf("hasher.Write: %v", err)
	}
	truncatedRefresh := hasher.Sum(nil)

	// bcrypt-хэш
	if err := bcrypt.CompareHashAndPassword(storedRefresh.Hash, []byte(truncatedRefresh)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", "", ErrWrongRefresh
		}
		return "", "", fmt.Errorf("bcrypt.CompareHashAndPassword: %v", err)
	}

	// валидация прошла успешно - создаем новую пару токенов
	return as.Create(ctx, id, ip)
}

func (as *AuthService) sendWarning(ctx context.Context, id, ip string) error {
	guid, err := guid.FromString(id)
	if err != nil {
		return fmt.Errorf("guid.FromString: %v", err)
	}

	mail, err := as.repo.GetEmail(ctx, guid)
	if err != nil {
		return err
	}

	return smtp.SendMail(
		fmt.Sprintf("%v:%v", as.conf.SMTP.Hostname, as.conf.SMTP.Port),
		as.auth,
		as.conf.SMTP.Username,
		[]string{mail},
		[]byte(fmt.Sprintf("Subject: Warning\nSomeone tried to accces your account from %v\n", ip)),
	)
}
