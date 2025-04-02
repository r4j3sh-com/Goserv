package main

import (
	"net/http"
	"time"

	"github.com/r4j3sh-com/Goserv/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	getRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing refresh token", err)
		return
	}

	// Validate the refresh token from db
	uDetails, err := cfg.db.GetRefreshTokenByToken(r.Context(), getRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}
	if uDetails.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked", nil)
		return
	}
	if uDetails.UserID.UUID.String() == " " || !uDetails.UserID.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID in refresh token", nil)
		return
	}
	newAccessToken, err := auth.MakeJWT(uDetails.UserID.UUID, cfg.tokenSecret, time.Hour*1)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new JWT", err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newAccessToken,
	})
}
