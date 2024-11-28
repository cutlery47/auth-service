package app

import (
	"context"
	"log"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/config"
	"github.com/cutlery47/auth-service/internal/repository"
	"github.com/cutlery47/auth-service/internal/service"
)

func RunAgent() {
	ctx := context.Background()

	conf, err := config.New()
	if err != nil {
		log.Fatal("error when reading config: ", err)
	}

	repo := repository.NewMock(conf.Repository)

	srv := service.NewAuthService(repo, conf.Service)

	// 1) Тест при правильных параметрах
	ip1 := "localhost"
	id1, _ := guid.NewV4()

	_, refresh1, err := srv.Create(ctx, id1, ip1)
	if err != nil {
		log.Fatal("error: ", err)
	}

	_, _, err = srv.Refresh(ctx, id1, ip1, refresh1)
	if err != nil {
		log.Fatal("error: ", err)
	}

	// 2) Тест на попытку рефрешнуть чужим токеном
	ip2 := "localhost"
	id2, _ := guid.NewV4()

	_, _, err = srv.Create(ctx, id2, ip2)
	if err != nil {
		log.Fatal("error: ", err)
	}

	// попытка рефрешнуть чужим
	_, _, err = srv.Refresh(ctx, id2, ip2, refresh1)
	if err != nil {
		log.Println("TEST WRONG REFRESH PASSED: ", err)
	}

	// 3) Тест на попытку рефрешнуть токен с другим ip
	ip3 := "localhost1"
	id3, _ := guid.NewV4()

	_, refresh3, err := srv.Create(ctx, id3, ip3)
	if err != nil {
		log.Fatal("error: ", err)
	}

	// попытка рефрешнуть другим ip
	_, _, err = srv.Refresh(ctx, id3, "localhost", refresh3)
	if err != nil {
		log.Println("TEST WRONG IP PASSED: ", err)
	}

	// 4) Тест на попытку рефрешнуть просроченный токен
	ip4 := "localhost"
	id4, _ := guid.NewV4()

	_, _, err = srv.Create(ctx, id4, ip4)
	if err != nil {
		log.Fatal("error: ", err)
	}

	// 5) Тест на попытку рефрешнуть рандомный набор символов
	_, _, err = srv.Refresh(ctx, id4, "1212313", "12312313")
	if err != nil {
		log.Println("TEST RANDOM REFRESH PASSED: ", err)
	}

}
