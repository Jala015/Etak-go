package helix

import (
	"testing"
)

// newTestStore points a Store at a test server.
func newTestStore(t *testing.T, baseURL string) *Store {
	t.Helper()
	s, err := New(baseURL)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	return s
}

func TestNew_invalidBaseURL(t *testing.T) {
	if _, err := New("not-a-url"); err == nil {
		t.Fatal("New() accepted a baseURL without scheme/host")
	}
}
