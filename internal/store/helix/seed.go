package helix

import (
	"context"
	"fmt"

	helixdb "github.com/helixdb/helix-db/sdks/go"

	"etak-go/internal/model"
)

// fotossinteseIslands is the seed for the protocol's example (section 3 of the old plan).
var fotossinteseIslands = []model.Island{
	{ID: "luz-solar", Name: "Luz solar", Definition: "A luz emitida pelo sol.", Status: model.StatusActive},
	{ID: "atp", Name: "ATP", Definition: "Molécula que transporta energia química.", Status: model.StatusActive},
	{ID: "clorofila", Name: "Clorofila", Definition: "Pigmento que absorve luz para a fotossíntese.", Status: model.StatusActive},
	{ID: "fotossintese", Name: "Fotossíntese", Definition: "Processo que converte luz em energia química.", Status: model.StatusActive},
}

// createIslandIDIndex builds the creation of the secondary equality index over
// Island.id. Idempotent.
func createIslandIDIndex() helixdb.Request {
	q := helixdb.WriteQuery("create_island_id_index")

	return q.VarAs("index", helixdb.G().
		CreateIndexIfNotExists(helixdb.NodeEqualityIndex("Island", "id")),
	).Returning("index")
}

// addIsland builds the insertion of an island as an Island node.
func addIsland(i model.Island) helixdb.Request {
	q := helixdb.WriteQuery("create_island")

	return q.VarAs("island", helixdb.G().
		AddN("Island", helixdb.Props{
			helixdb.Prop("id", q.ParamString("id", i.ID)),
			helixdb.Prop("name", q.ParamString("name", i.Name)),
			helixdb.Prop("definition", q.ParamString("definition", i.Definition)),
			helixdb.Prop("status", q.ParamString("status", string(i.Status))),
		}),
	).Returning("island")
}

// Bootstrap creates the necessary indexes (idempotent).
func (s *Store) Bootstrap(ctx context.Context) error {
	if err := s.client.Exec(ctx, createIslandIDIndex(), nil); err != nil {
		return fmt.Errorf("create Island.id index: %w", err)
	}
	return nil
}

// SeedDefault runs Bootstrap and inserts the fotossíntese seed.
func (s *Store) SeedDefault(ctx context.Context) error {
	if err := s.Bootstrap(ctx); err != nil {
		return err
	}
	return s.SeedIslands(ctx, fotossinteseIslands...)
}

// SeedIslands inserts each island as an Island node (one add_n per island).
func (s *Store) SeedIslands(ctx context.Context, islands ...model.Island) error {
	for _, i := range islands {
		if err := s.client.Exec(ctx, addIsland(i), nil); err != nil {
			return fmt.Errorf("insert island %q: %w", i.ID, err)
		}
	}
	return nil
}
