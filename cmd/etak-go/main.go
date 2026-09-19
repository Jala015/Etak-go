package main

import (
	"log"
	"net/http"
	"os"

	"etak-go/internal/api"
	"etak-go/internal/service"
	"etak-go/internal/store/helix"
)

func main() {
	helixAddr := os.Getenv("HELIX_ADDR")
	if helixAddr == "" {
		helixAddr = "http://localhost:6969"
	}

	st, err := helix.New(helixAddr) // Real store (HelixDB via POST /v2/query)
	if err != nil {
		log.Fatal(err)
	}
	svc := service.New(st)  // protocol logic
	handler := api.New(svc) // http.Handler

	log.Printf("etak-go listening at :8080 (HelixDB at %s)", helixAddr)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
