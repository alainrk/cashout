package repository

import (
	"cashout/internal/model"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type Auth struct {
	Repository
}

// GenerateAuthToken generates a random alphanumeric token
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// GenerateSessionID generates a secure session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateAuthToken creates a new auth token for a user
func (r *Auth) CreateAuthToken(tgID int64) (*model.AuthToken, error) {
	token, err := generateRandomToken(6) // 6-character alphanumeric code
	if err != nil {
		return nil, err
	}

	authToken := &model.AuthToken{
		TgID:      tgID,
		Token:     token,
		Status:    model.AuthStatusPending,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5 minutes expiry
	}

	if err := r.DB.CreateAuthToken(authToken); err != nil {
		return nil, err
	}

	return authToken, nil
}

// GetAuthToken retrieves an auth token by token string
func (r *Auth) GetAuthToken(token string) (*model.AuthToken, error) {
	return r.DB.GetAuthToken(token)
}

// VerifyAuthToken marks a token as verified and returns the user
func (r *Auth) VerifyAuthToken(token string) (*model.User, error) {
	authToken, err := r.DB.GetAuthToken(token)
	if err != nil {
		return nil, err
	}

	if !authToken.IsValid() {
		return nil, ErrInvalidToken
	}

	// Mark token as verified
	authToken.Status = model.AuthStatusVerified
	if err := r.DB.UpdateAuthToken(authToken); err != nil {
		return nil, err
	}

	// Get user
	user, err := r.DB.GetUser(authToken.TgID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateWebSession creates a new web session for a user
func (r *Auth) CreateWebSession(tgID int64) (*model.WebSession, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &model.WebSession{
		ID:        sessionID,
		TgID:      tgID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours session
	}

	if err := r.DB.CreateWebSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetWebSession retrieves a web session by ID
func (r *Auth) GetWebSession(sessionID string) (*model.WebSession, error) {
	return r.DB.GetWebSession(sessionID)
}

// DeleteWebSession deletes a web session
func (r *Auth) DeleteWebSession(sessionID string) error {
	return r.DB.DeleteWebSession(sessionID)
}

// CleanupExpiredTokens removes expired auth tokens and sessions
func (r *Auth) CleanupExpiredTokens() error {
	return r.DB.CleanupExpiredAuthData()
}
