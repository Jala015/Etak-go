// Command etak-seed creates the indexes and inserts the example seed (photosynthesis)
// into a local HelixDB instance. Re-runnable (local storage is in memory).

package main

import (
	"context"
	"log"
	"os"

	"etak-go/internal/store/helix"
)

func main() {
	addr := os.Getenv("HELIX_ADDR")
	if addr == "" {
		addr = "http://localhost:6969"
	}

	st, err := helix.New(addr)
	if err != nil {
		log.Fatal(err)
	}
	if err := st.SeedDefault(context.Background()); err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Printf("seed applied on %s", addr)
}
