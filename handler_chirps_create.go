package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/r4j3sh-com/Goserv/internal/auth"
	"github.com/r4j3sh-com/Goserv/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		//UserID uuid.NullUUID `json:"user_id"`
	}

	/* token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", nil)
		return
	} */

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	/* if !params.UserID.Valid {
		respondWithError(w, http.StatusBadRequest, "User ID is required", nil)
		return
	} */

	/* const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedWords := []string{}

	for _, word := range strings.Fields(params.Body) {
		lowercaseWord := strings.ToLower(word)
		replaced := false
		for _, badWord := range badWords {
			if lowercaseWord == badWord {
				cleanedWords = append(cleanedWords, "****")
				replaced = true
				break
			}
		}
		if !replaced {
			cleanedWords = append(cleanedWords, word)
		}
	}

	cleanedBody := strings.Join(cleanedWords, " ") */

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error validating chirp", err)
		return
	}

	//userId := params.UserID
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}
	userIdFromToken, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}
	userId := uuid.NullUUID{
		Valid: true,
		UUID:  userIdFromToken,
	}

	// add records to the database
	chirp, err := cfg.db.AddNewChrip(r.Context(), database.AddNewChripParams{
		Body:   cleanedBody,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding chirp to the database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    (chirp.UserID).UUID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

/* func (cfg *apiConfig) handlerChipsGet(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChrips(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps from the database", err)
		return
	}
	allChrips := make([]Chrip, len(chirps))
	for i := range chirps {
		allChrips[i] = Chrip{
			ID:        chirps[i].ID,
			CreatedAt: chirps[i].CreatedAt,
			UpdatedAt: chirps[i].UpdatedAt,
			Body:      chirps[i].Body,
			UserID:    (chirps[i].UserID).UUID,
		}
	}
	respondWithJSON(w, http.StatusOK, allChrips)
} */
