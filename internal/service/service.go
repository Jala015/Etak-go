package service

import (
	"context"
	"errors"

	"etak-go/internal/model"
	"etak-go/internal/store"
)

// ErrNotFound translates the store's ErrNotFound to the API's 404.
// We keep both: the service is the boundary between the store and HTTP.
var ErrNotFound = errors.New("not found")

type Service struct {
	store store.Store
}

func New(s store.Store) *Service {
	return &Service{store: s}
}

// GetIsland returns the island by id. If the id points to a ghost
// island (status deprecated with superseded_by), it follows the chain until the
// active entity — the recursive resolution the protocol requires.
func (s *Service) GetIsland(ctx context.Context, id string) (*model.Island, error) {
	seen := map[string]bool{} // tracks already visited islands.

	for {
		island, err := s.store.GetIsland(ctx, id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return nil, ErrNotFound // mapped from store; only the service's error reaches the API
			}
			return nil, err // pass the error through unchanged
		}

		// Switch on status so future statuses with different handling
		// can be added as separate cases (current unknown statuses
		// fall through to default as a direct response).
		switch island.Status {
		case model.StatusDeprecated:
			// Deprecated without superseded_by: treat as alive (direct response).
			if island.SupersededBy == nil {
				return island, nil
			}
			// Cycle guard: each id must appear at most once in the chain.
			if seen[id] {
				return nil, errors.New("cycle detected on superseded_by")
			}
			// Deprecated with superseded_by: ghost — follow one chain link.
			seen[id] = true
			id = *island.SupersededBy
			continue

		default: // active, draft, or any future status: direct response.
			return island, nil
		}

	}
}
