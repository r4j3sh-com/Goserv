package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckJWT(t *testing.T) {
	userId := uuid.New()
	userId2 := uuid.New()
	tokenSecret := "testSecret"

	token1, _ := MakeJWT(userId, tokenSecret, time.Minute*5)
	token2, _ := MakeJWT(userId2, tokenSecret, time.Second*1)

	tests := []struct {
		name        string
		token       string
		tokenSecret string
		wantUserId  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			token:       token1,
			tokenSecret: tokenSecret,
			wantUserId:  userId,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			token:       "invalidToken",
			tokenSecret: tokenSecret,
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Expired token",
			token:       token2,
			tokenSecret: tokenSecret,
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "No token provided",
			token:       "",
			tokenSecret: tokenSecret,
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Invalid token secret",
			token:       token1,
			tokenSecret: "invalidSecret",
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Expired token" {
				time.Sleep(time.Second * 2)
			}
			gotUserId, err := ValidateJWT(tt.token, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotUserId != tt.wantUserId {
				t.Errorf("ValidateJWT() userId = %v, want %v", gotUserId, tt.wantUserId)
			}
		})
	}
}
