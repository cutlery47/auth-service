package app

import (
	"context"
	"log"
	"os"

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

	_ = service.NewAuthService(repo, conf.Service)
}
