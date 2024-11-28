package repository

import (
	"sync"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/cutlery47/auth-service/internal/models"
	"github.com/google/uuid"
)

type MockRepository struct {
	tokens []models.OutRefresh
	mu     *sync.RWMutex
}

func NewMock() *MockRepository {
	return &MockRepository{
		tokens: []models.OutRefresh{},
		mu:     &sync.RWMutex{},
	}
}

func (mr *MockRepository) Create(refresh models.InRefresh) error {
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

func (mr *MockRepository) Get(id guid.GUID) (models.OutRefresh, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	for _, el := range mr.tokens {
		if el.UserId == id {
			return el, nil
		}
	}

	return models.OutRefresh{}, ErrNotFound
}
