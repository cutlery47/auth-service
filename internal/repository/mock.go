package repository

import (
	"context"
	"sync"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/config"
	"github.com/cutlery47/auth-service/internal/models"
	"github.com/google/uuid"
)

type MockRepository struct {
	tokens []models.OutRefresh
	mu     *sync.RWMutex

	conf config.Repository
}

func NewMock(conf config.Repository) *MockRepository {
	return &MockRepository{
		tokens: []models.OutRefresh{},
		mu:     &sync.RWMutex{},

		conf: conf,
	}
}

func (mr *MockRepository) Create(ctx context.Context, refresh models.InRefresh) error {
	id := uuid.New()

	entry := models.OutRefresh{
		Id:        id,
		InRefresh: refresh,
	}

	mr.mu.Lock()
	mr.tokens = append(mr.tokens, entry)
	mr.mu.Unlock()

	return nil
}

func (mr *MockRepository) Get(ctx context.Context, id guid.GUID) (models.OutRefresh, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	for _, el := range mr.tokens {
		if el.UserId == id {
			return el, nil
		}
	}

	return models.OutRefresh{}, ErrNotFound
}

func (mr *MockRepository) GetEmail(ctx context.Context, id guid.GUID) (string, error) {
	return mr.conf.Receiver, nil
}
