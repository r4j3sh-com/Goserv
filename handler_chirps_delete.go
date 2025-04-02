package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/r4j3sh-com/Goserv/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	getAccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing access token", err)
		return
	}
	verifyToken, err := auth.ValidateJWT(getAccessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired access token", err)
		return
	}

	chripsID := r.PathValue("chirpID")
	if chripsID == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required", nil)
		return
	}
	parsedChirpID, err := uuid.Parse(chripsID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirps ID", err)
		return
	}

	getChripDetails, err := cfg.db.GetChripByID(r.Context(), parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if getChripDetails.UserID.UUID != verifyToken {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp", nil)
		return
	}

	err = cfg.db.DeleteChripsByID(r.Context(), parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp from the database", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
