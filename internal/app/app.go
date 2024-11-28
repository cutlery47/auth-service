package app

import (
	"context"
	"log"
	"os"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/config"
	"github.com/cutlery47/auth-service/internal/repository"
	"github.com/cutlery47/auth-service/internal/service"
	"github.com/sirupsen/logrus"
)

func Run() {
	ctx := context.Background()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	conf, err := config.New()
	if err != nil {
		log.Fatal("error when reading config: ", err)
	}

	repo, err := repository.NewAuthRepository(ctx, conf.Repository)
	if err != nil {
		log.Fatal("error when initializing app: ", err)
	}

	srv := service.NewAuthService(repo, conf.Service)

	id, _ := guid.NewV4()
	ip := "localhost"

	_, refresh, err := srv.Create(ctx, id, ip)
	if err != nil {
		log.Fatal("error: ", err)
	}

	_, _, err = srv.Refresh(ctx, id, "localhost", refresh)
	if err != nil {
		log.Fatal("error: ", err)
	}

}
