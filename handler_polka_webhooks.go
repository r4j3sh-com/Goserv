package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/r4j3sh-com/Goserv/internal/auth"
	"github.com/r4j3sh-com/Goserv/internal/database"
)

type PolkaWebhookRequest struct {
	Event string
	Data  map[string]interface{}
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing API key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var polkaWebhookRequest PolkaWebhookRequest
	err = decoder.Decode(&polkaWebhookRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Polka Webhook request", err)
		return
	}

	if polkaWebhookRequest.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(polkaWebhookRequest.Data["user_id"].(string))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving user details", err)
		return
	}

	userDetails, err := cfg.db.GetUsersByID(r.Context(), userId) // Extract user details from the request data
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving user details", err)
		return
	}

	_, err = cfg.db.UpgradeUser(r.Context(), database.UpgradeUserParams{
		IsChirpyRed: true,
		ID:          userDetails.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
