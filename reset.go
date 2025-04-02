package main

import (
	"net/http"
	"strconv"
)

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if c.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset endpoint is only available in development environment", nil)
		return
	}
	c.fileserveHits.Store(0)
	resetDb, err := c.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset Done"))
	w.Write([]byte("\nTotal Users deleted: " + strconv.Itoa(len(resetDb)) + "\n"))
}
