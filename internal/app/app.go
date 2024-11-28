package app

import (
	"context"
	"fmt"
	"log"

	"github.com/cutlery47/auth-service/internal/config"
	v1 "github.com/cutlery47/auth-service/internal/controller/http/v1"
	"github.com/cutlery47/auth-service/internal/repository"
	"github.com/cutlery47/auth-service/internal/service"
	"github.com/cutlery47/auth-service/internal/utils"
	"github.com/cutlery47/auth-service/pkg/httpserver"
	"github.com/cutlery47/auth-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	_ "github.com/cutlery47/auth-service/docs"
)

//	@title			Authentication Service
//	@version		0.0.1
//	@description	This is an authentication service

//	@contact.name	Ivanchenko Arkhip
//	@contact.email	kitchen_cutlery@mail.ru

//	@BasePath	/

func Run() error {
	ctx := context.Background()

	logrus.Debug("reading config...")
	conf, err := config.New()
	if err != nil {
		log.Fatal("error when reading config: ", err)
	}

	logrus.Debug("creating loggers...")
	infoFd, err := utils.CreateAndOpen(conf.Logger.InfoPath)
	if err != nil {
		return fmt.Errorf("error when creating info log file: %v", err)
	}

	errFd, err := utils.CreateAndOpen(conf.Logger.ErrorPath)
	if err != nil {
		return fmt.Errorf("error when creating error log file: %v", err)
	}

	infoLog := logger.WithFormat(logger.WithFile(logger.New(logrus.InfoLevel), infoFd), &logrus.JSONFormatter{})
	errLog := logger.WithFormat(logger.WithFile(logger.New(logrus.ErrorLevel), errFd), &logrus.JSONFormatter{})

	logrus.Debug("initializing repository...")
	repo, err := repository.NewAuthRepository(ctx, conf.Repository)
	if err != nil {
		log.Fatal("error when initializing app: ", err)
	}

	logrus.Debug("initializing service...")
	srv := service.NewAuthService(repo, conf.Service)

	logrus.Debug("initializing controller...")
	echo := echo.New()
	v1.NewController(echo, srv, infoLog, errLog)

	logrus.Debug("initializing http server...")
	server := httpserver.New(echo, conf.HTTPServer)

	return server.Run(ctx)
}
