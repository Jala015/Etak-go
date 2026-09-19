package api

import (
	"net/http"

	"etak-go/internal/service"
)

// API holds dependencies and loads handlers.
type API struct {
	svc *service.Service
}

// New mounts the router. The ServeMux of Go 1.22+ understands method + path.
func New(svc *service.Service) http.Handler {
	a := &API{svc: svc}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", a.healthz)
	mux.HandleFunc("GET /v1/islands/{id}", a.getIsland) // wildcard {id}
	// Fase 4+: POST /v1/islands, PUT /v1/islands/{id}, DELETE /v1/islands/{id}...

	return mux
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
