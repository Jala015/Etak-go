package api

import (
	"errors"
	"net/http"

	"etak-go/internal/service"
)

func (a *API) getIsland(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // extracts {id} from route

	island, err := a.svc.GetIsland(r.Context(), id)
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, island)
}
