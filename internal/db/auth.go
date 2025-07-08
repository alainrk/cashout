package db

import (
	"cashout/internal/model"
	"time"
)

// CreateAuthToken creates a new auth token
func (db *DB) CreateAuthToken(token *model.AuthToken) error {
	return db.conn.Create(token).Error
}

// GetAuthToken retrieves an auth token by token string
func (db *DB) GetAuthToken(token string) (*model.AuthToken, error) {
	var authToken model.AuthToken
	result := db.conn.Where("token = ?", token).First(&authToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &authToken, nil
}

// UpdateAuthToken updates an auth token
func (db *DB) UpdateAuthToken(token *model.AuthToken) error {
	return db.conn.Save(token).Error
}

// CreateWebSession creates a new web session
func (db *DB) CreateWebSession(session *model.WebSession) error {
	return db.conn.Create(session).Error
}

// GetWebSession retrieves a web session by ID
func (db *DB) GetWebSession(sessionID string) (*model.WebSession, error) {
	var session model.WebSession
	result := db.conn.Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

// GetWebSessionWithUser retrieves a web session with user data
func (db *DB) GetWebSessionWithUser(sessionID string) (*model.WebSession, error) {
	var session model.WebSession
	result := db.conn.Preload("User").Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

// DeleteWebSession deletes a web session
func (db *DB) DeleteWebSession(sessionID string) error {
	return db.conn.Delete(&model.WebSession{}, "id = ?", sessionID).Error
}

// CleanupExpiredAuthData removes expired auth tokens and sessions
func (db *DB) CleanupExpiredAuthData() error {
	now := time.Now()

	// Delete expired auth tokens
	if err := db.conn.Where("expires_at < ?", now).Delete(&model.AuthToken{}).Error; err != nil {
		return err
	}

	// Delete expired web sessions
	if err := db.conn.Where("expires_at < ?", now).Delete(&model.WebSession{}).Error; err != nil {
		return err
	}

	return nil
}
