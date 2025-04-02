package main

import (
	"net/http"

	"github.com/r4j3sh-com/Goserv/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	getRefreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing refresh token", err)
		return
	}

	if getRefreshToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing refresh token", nil)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), getRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
