package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("author_id") != "" {
		cfg.handlerChirpsRetrieveByAuthor(w, r)
		return
	}
	dbChirps, err := cfg.db.GetChrips(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    (dbChirp.UserID).UUID,
			Body:      dbChirp.Body,
		})
	}

	// New sorting logic
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "" {
		switch sortOrder {
		case "asc":
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			})
		case "desc":
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
		default:
			respondWithError(w, http.StatusBadRequest, "Invalid sort order. Use 'asc' or 'desc'", nil)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChipRetriveByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirps ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChripByID(r.Context(), parsedId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	if dbChirp.Body == "" {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    (dbChirp.UserID).UUID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieveByAuthor(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")

	parsedId, err := uuid.Parse(authorID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID", err)
		return
	}
	authorNullUUID := uuid.NullUUID{UUID: parsedId, Valid: true}
	retriveChirps, err := cfg.db.GetChirpsByAuthorID(r.Context(), authorNullUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range retriveChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    (dbChirp.UserID).UUID,
			Body:      dbChirp.Body,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
