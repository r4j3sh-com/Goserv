package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/r4j3sh-com/Goserv/internal/auth"
	"github.com/r4j3sh-com/Goserv/internal/database"
)

type UserDetails struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}

type parameters struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (c *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	// IF expires_in_seconds is provided, generate JWT with expiration time
	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, c.tokenSecret, expirationTime)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	newUser := UserDetails{
		ID:          user.ID.String(),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		Email:       user.Email,
		Token:       token,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}

func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}
	user, err := c.db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", nil)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}
	// IF expires_in_seconds is provided, generate JWT with expiration time
	var expiresIn int
	if params.ExpiresInSeconds != 0 {
		expiresIn = params.ExpiresInSeconds
	} else {
		expiresIn = 3600 // 1 hour by default if expires_in_seconds is not provided
	}

	token, err := auth.MakeJWT(user.ID, c.tokenSecret, time.Second*time.Duration(expiresIn))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	// Create New Refresh Token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	refreshTokenExpiresAt := time.Now().Add(time.Duration(24*60) * time.Hour)

	savedrefreshToken, err := c.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: uuid.NullUUID{
			Valid: true,
			UUID:  user.ID,
		},
		ExpiresAt: refreshTokenExpiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	user_details := UserDetails{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		Email:        user.Email,
		Token:        token,
		RefreshToken: savedrefreshToken.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, user_details)
}

func (c *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	getAccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing access token", err)
		return
	}
	verifyToken, err := auth.ValidateJWT(getAccessToken, c.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired access token", err)
		return
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Email == "" && params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	var hashedPassword string
	hashedPassword, err = auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	user, err := c.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             verifyToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	user_details := UserDetails{
		ID:          user.ID.String(),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		Email:       user.Email,
		Token:       getAccessToken,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, user_details)
}
