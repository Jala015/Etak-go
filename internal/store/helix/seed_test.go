package helix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"etak-go/internal/model"
)

func TestBootstrap_createsSecondaryIndex(t *testing.T) {
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(body, &m)
		bodies = append(bodies, m)
		io.WriteString(w, `{"index":[]}`)
	}))
	defer srv.Close()

	if err := newTestStore(t, srv.URL).Bootstrap(context.Background()); err != nil {
		t.Fatalf("Bootstrap() error: %v", err)
	}
	if len(bodies) == 0 {
		t.Fatal("Bootstrap sent no request")
	}
	if bodies[0]["request_type"] != "write" {
		t.Fatalf("first request is not write: %#v", bodies[0])
	}
}

func TestSeedIslands_oneAddNPerIsland(t *testing.T) {
	var count int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		io.WriteString(w, `{"island":[{"$id":0}]}`)
	}))
	defer srv.Close()

	err := newTestStore(t, srv.URL).SeedIslands(context.Background(),
		model.Island{ID: "a", Name: "A", Status: model.StatusActive},
		model.Island{ID: "b", Name: "B", Status: model.StatusActive},
	)
	if err != nil {
		t.Fatalf("SeedIslands() error: %v", err)
	}
	if count != 2 {
		t.Fatalf("sent %d requests, want 2", count)
	}
}
