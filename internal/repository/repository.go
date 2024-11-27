package repository

import (
	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/models"
)

type Repository interface {
	Get(id guid.GUID) (models.OutRefresh, error)
	Create(refresh models.InRefresh) error
}
