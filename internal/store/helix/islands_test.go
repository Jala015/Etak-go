package helix

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"etak-go/internal/model"
	"etak-go/internal/store"
)

func TestGetIsland_found(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/query" || r.Method != http.MethodPost {
			t.Errorf("request = %s %s, want POST /v2/query", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"island":[{"$id":0,"id":"clorofila","name":"Clorofila","definition":"Pigmento.","status":"active"}]}`)
	}))
	defer srv.Close()

	got, err := newTestStore(t, srv.URL).GetIsland(context.Background(), "clorofila")
	if err != nil {
		t.Fatalf("GetIsland() error: %v", err)
	}
	if got.ID != "clorofila" || got.Name != "Clorofila" || got.Status != "active" {
		t.Fatalf("GetIsland() = %+v", got)
	}
}

func TestGetIsland_notFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"island":[]}`)
	}))
	defer srv.Close()

	_, err := newTestStore(t, srv.URL).GetIsland(context.Background(), "inexistente")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("error = %v, want store.ErrNotFound", err)
	}
}

func TestIslandRowToModel_full(t *testing.T) {
	superseded, derived := "nova", "origem"
	got := islandRow{
		ID:           "clorofila",
		Name:         "Clorofila",
		Definition:   "Pigmento.",
		Status:       "active",
		SupersededBy: &superseded,
		DerivedFrom:  &derived,
	}.toModel()

	if got.ID != "clorofila" || got.Name != "Clorofila" || got.Definition != "Pigmento." {
		t.Fatalf("toModel() = %+v", got)
	}
	if got.Status != model.StatusActive {
		t.Fatalf("Status = %q, want active", got.Status)
	}
	if got.SupersededBy == nil || *got.SupersededBy != "nova" {
		t.Fatalf("SupersededBy = %v", got.SupersededBy)
	}
	if got.DerivedFrom == nil || *got.DerivedFrom != "origem" {
		t.Fatalf("DerivedFrom = %v", got.DerivedFrom)
	}
}

func TestIslandRowToModel_optionalAbsentOrEmptyBecomeNil(t *testing.T) {
	empty := ""
	got := islandRow{ID: "a", Status: "active", SupersededBy: &empty}.toModel()

	if got.SupersededBy != nil {
		t.Fatalf("SupersededBy = %v, want nil for empty string", got.SupersededBy)
	}
	if got.DerivedFrom != nil {
		t.Fatalf("DerivedFrom = %v, want nil when absent", got.DerivedFrom)
	}
}

func TestIslandRowToModel_unknownStatusBecomesDraft(t *testing.T) {
	got := islandRow{ID: "a", Status: "inventado"}.toModel()
	if got.Status != model.StatusDraft {
		t.Fatalf("Status = %q, want draft", got.Status)
	}
}
