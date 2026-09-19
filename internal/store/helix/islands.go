package helix

import (
	"context"
	"fmt"

	helixdb "github.com/helixdb/helix-db/sdks/go"

	"etak-go/internal/model"
	"etak-go/internal/store"
)

// islandProps are the properties we project when reading an Island.
// "$id" is the internal node id in HelixDB; "id" is the canonical protocol key.
var islandProps = []string{"$id", "id", "name", "definition", "status", "superseded_by", "derived_from"}

// islandByID builds the read of an island by canonical id. The label filter
// ensures only Island nodes enter the result.
func islandByID(id string) helixdb.Request {
	q := helixdb.ReadQuery("island_by_id")
	param := q.ParamString("id", id)

	return q.VarAs("island", helixdb.G().
		NWithLabel("Island").
		Where(helixdb.PredEq("id", param)).
		ValueMap(islandProps...),
	).Returning("island")
}

// GetIsland returns the island by canonical id. Empty result → ErrNotFound.
func (s *Store) GetIsland(ctx context.Context, id string) (*model.Island, error) {
	var resp struct {
		Island []islandRow `json:"island"`
	}
	if err := s.client.Exec(ctx, islandByID(id), &resp); err != nil {
		return nil, fmt.Errorf("fetch island %q: %w", id, err)
	}
	if len(resp.Island) == 0 {
		return nil, store.ErrNotFound
	}
	return resp.Island[0].toModel(), nil
}

// islandRow is the row projected by value_map for an Island node. Missing properties become zero values on unmarshal.
type islandRow struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Definition   string  `json:"definition"`
	Status       string  `json:"status"`
	SupersededBy *string `json:"superseded_by"`
	DerivedFrom  *string `json:"derived_from"`
}

// toModel converts the row to the domain model. Missing or empty optional fields become nil.
func (r islandRow) toModel() *model.Island {
	return &model.Island{
		ID:           r.ID,
		Name:         r.Name,
		Definition:   r.Definition,
		Status:       model.ParseStatus(r.Status),
		SupersededBy: nilIfEmpty(r.SupersededBy),
		DerivedFrom:  nilIfEmpty(r.DerivedFrom),
	}
}
