package fake

import (
	"context"

	"etak-go/internal/model"
	"etak-go/internal/store"
)

// Store is an in-memory implementation of store.Store, used to test
// the service without needing HelixDB.
type Store struct {
	islands map[string]model.Island
}

func New() *Store {
	return &Store{islands: make(map[string]model.Island)}
}

// Seed populates the store with islands to set up test scenarios.
func (s *Store) Seed(islands ...model.Island) {
	for _, i := range islands {
		s.islands[i.ID] = i
	}
}

func (s *Store) GetIsland(ctx context.Context, id string) (*model.Island, error) {
	i, ok := s.islands[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	cp := i // return a copy, so the caller can't mutate the store
	return &cp, nil
}
