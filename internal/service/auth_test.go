package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/config"
	"github.com/cutlery47/auth-service/internal/repository"
)

var conf = config.Config{
	Service: config.Service{
		AccessTTL:  3 * time.Second,
		RefreshTTL: 3 * time.Second,
		Secret:     "somesecret",
		Cost:       10,
	},
	Repository: config.Repository{
		Receiver: "example@gmail.com",
	},
}

var srv *AuthService

func setup() {
	repo := repository.NewMock(conf.Repository)
	srv = NewAuthService(repo, conf.Service)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

// тест при корректных данных
func TestCorrect(t *testing.T) {
	ctx := context.Background()

	ip := "localhost"
	id, _ := guid.NewV4()

	_, refresh, err := srv.Create(ctx, id, ip)
	if err != nil {
		t.Fatal("error: ", err)
	}

	_, _, err = srv.Refresh(ctx, id, ip, refresh)
	if err != nil {
		t.Fatal("error: ", err)
	}
}

// тест попытки рефрешнуть свои токены чужими
func TestWrongToken(t *testing.T) {
	ctx := context.Background()

	var (
		ip1 = "localhost"
		ip2 = "localhost"

		id1, _ = guid.NewV4()
		id2, _ = guid.NewV4()
	)

	_, _, err := srv.Create(ctx, id1, ip1)
	if err != nil {
		t.Fatal("error: ", err)
	}

	_, refresh2, err := srv.Create(ctx, id2, ip2)
	if err != nil {
		t.Fatal("error: ", err)
	}

	// попытка рефрешнуть чужим
	_, _, err = srv.Refresh(ctx, id1, ip1, refresh2)
	if err != nil {
		if !errors.Is(err, ErrWrongRefresh) {
			t.Fatal("error: ", err)
		}
	}
}

func TestMalformedToken(t *testing.T) {
	ctx := context.Background()

	var (
		ip    = "localhost"
		id, _ = guid.NewV4()

		wrongRefresh = "2281337"
	)

	_, _, err := srv.Create(ctx, id, ip)
	if err != nil {
		t.Fatal("error: ", err)
	}

	_, _, err = srv.Refresh(ctx, id, ip, wrongRefresh)
	if err != nil {
		if !errors.Is(err, ErrMalformedToken) {
			t.Fatal("error: ", err)
		}
	}
}

func TestUnexistantId(t *testing.T) {
	ctx := context.Background()

	var (
		ip = "localhost"

		id1, _ = guid.NewV4()
		id2, _ = guid.NewV4()
	)

	_, refresh, err := srv.Create(ctx, id1, ip)
	if err != nil {
		t.Fatal("error: ", err)
	}

	// прокидываем id2 вместо id1
	_, _, err = srv.Refresh(ctx, id2, ip, refresh)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			t.Fatal("error: ", err)
		}
	}
}

func TestWrongId(t *testing.T) {
	ctx := context.Background()

	var (
		ip1 = "localhost"
		ip2 = "localhost"

		id1, _ = guid.NewV4()
		id2, _ = guid.NewV4()
	)

	_, refresh1, err := srv.Create(ctx, id1, ip1)
	if err != nil {
		t.Fatal("error: ", err)
	}

	_, _, err = srv.Create(ctx, id2, ip2)
	if err != nil {
		t.Fatal("error: ", err)
	}

	// прокидываем id2 вместо id1
	_, _, err = srv.Refresh(ctx, id2, ip1, refresh1)
	if err != nil {
		if !errors.Is(err, ErrWrongRefresh) {
			t.Fatal("error: ", err)
		}
	}
}

func TestWrongIp(t *testing.T) {
	ctx := context.Background()

	var (
		ip1    = "localhost"
		id1, _ = guid.NewV4()

		ip2 = "globalhost"
	)

	_, refresh, err := srv.Create(ctx, id1, ip1)
	if err != nil {
		t.Fatal("error: ", err)
	}

	// прокидываем ip2 место ip1
	_, _, err = srv.Refresh(ctx, id1, ip2, refresh)
	if err != nil {
		if !errors.Is(err, ErrWrongIp) {
			t.Fatal("error: ", err)
		}
	}
}
