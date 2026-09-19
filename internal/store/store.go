package store

import (
	"context"
	"errors"

	"etak-go/internal/model"
)

// ErrNotFound is a sentinel error: the service uses errors.Is to decide whether to respond with 404.
var ErrNotFound = errors.New("not found")

// Store isolates all communication with HelixDB. The service only knows this
// interface; tests use store/fake, production uses store/helix.
//
// The interface grows with each phase:
//   - Phase 4+: CreateIsland, UpdateIsland, DeleteIsland
//   - Phase 5+: FlowsFrom, FlowsTo, AllFlows (graph traversal and queries)
type Store interface {
	GetIsland(ctx context.Context, id string) (*model.Island, error)
}
