package helix

import (
	"fmt"
	"net/http"
	"time"

	helixdb "github.com/helixdb/helix-db/sdks/go"

	"etak-go/internal/store"
)

// guarantees at compile time that *Store implements store.Store.
var _ store.Store = (*Store)(nil)

// requestTimeout limits each call to HelixDB. The SDK uses http.DefaultClient,
// which has no timeout, so we inject our own client.
const requestTimeout = 30 * time.Second

// Store implements store.Store over a HelixDB instance, speaking the
// POST /v2/query endpoint through the official SDK.
type Store struct {
	client *helixdb.Client
}

// New creates the Store pointing to baseURL (e.g., "http://localhost:6969").
func New(baseURL string) (*Store, error) {
	client, err := helixdb.NewClient(baseURL,
		helixdb.WithHTTPClient(&http.Client{Timeout: requestTimeout}),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to HelixDB at %q: %w", baseURL, err)
	}
	return &Store{client: client}, nil
}
