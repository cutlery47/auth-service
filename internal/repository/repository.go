package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/config"
	"github.com/cutlery47/auth-service/internal/models"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

type Repository interface {
	Get(ctx context.Context, id guid.GUID) (models.OutRefresh, error)
	GetEmail(ctx context.Context, id guid.GUID) (string, error)
	Create(ctx context.Context, refresh models.InRefresh) error
}

type AuthRepository struct {
	db *sql.DB

	conf config.Repository
}

func NewAuthRepository(ctx context.Context, conf config.Repository) (*AuthRepository, error) {
	url := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		conf.Postgres.Username,
		conf.Postgres.Password,
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.DB,
	)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// тайм-аут для подключения к бд
	timeoutCtx, cancel := context.WithTimeout(ctx, conf.Postgres.Timeout)
	defer cancel()

	// пингуем бд, чтобы проверить, что она запущена и принимает соединения
	err = db.PingContext(timeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish connection with postgres: %v", err)
	}
	logrus.Debug("successfully established postgres connection!")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres.WithInstance: %v", err)
	}

	migrations := fmt.Sprintf("file://%v", conf.Postgres.Migrations)
	m, err := migrate.NewWithDatabaseInstance(migrations, conf.Postgres.DB, driver)
	if err != nil {
		return nil, fmt.Errorf("migrate.NewWithDatabaseInstance: %v", err)
	}

	// мигрируемся
	logrus.Debug("applying migrations...")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Debug("nothing to migrate")
		} else {
			return nil, fmt.Errorf("error when migrating: %v", err)
		}
	} else {
		logrus.Debug("migrated successfully!")
	}
	defer m.Close()

	return &AuthRepository{
		db:   db,
		conf: conf,
	}, nil
}

func (ar *AuthRepository) Get(ctx context.Context, id guid.GUID) (models.OutRefresh, error) {
	return models.OutRefresh{}, ErrNotImplemented
}

func (ar *AuthRepository) Create(ctx context.Context, refresh models.InRefresh) error {
	return ErrNotImplemented
}

// "типо" получаем почту пользователя из бд и отправляем в сервис
func (ar *AuthRepository) GetEmail(ctx context.Context, id guid.GUID) (string, error) {
	return ar.conf.Receiver, nil
}
