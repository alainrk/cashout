package db

import (
	"cashout/internal/model"
	"time"

	"gorm.io/gorm"
)

// WebAuthn Credential Operations

// CreateWebAuthnCredential creates a new WebAuthn credential
func (db *DB) CreateWebAuthnCredential(cred *model.WebAuthnCredential) error {
	return db.conn.Create(cred).Error
}

// CreateWebAuthnCredentialWithSessionCleanup creates a credential and deletes the session in a transaction
func (db *DB) CreateWebAuthnCredentialWithSessionCleanup(cred *model.WebAuthnCredential, sessionID string) error {
	return db.conn.Transaction(func(tx *gorm.DB) error {
		// Create credential
		if err := tx.Create(cred).Error; err != nil {
			return err
		}

		// Cleanup session (non-critical, don't fail transaction)
		_ = tx.Delete(&model.WebAuthnSession{}, "id = ?", sessionID).Error

		return nil
	})
}

// GetWebAuthnCredential retrieves a credential by its ID
func (db *DB) GetWebAuthnCredential(credentialID []byte) (*model.WebAuthnCredential, error) {
	var cred model.WebAuthnCredential
	result := db.conn.Where("id = ?", credentialID).First(&cred)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cred, nil
}

// GetWebAuthnCredentialsByUser retrieves all credentials for a user
func (db *DB) GetWebAuthnCredentialsByUser(tgID int64) ([]model.WebAuthnCredential, error) {
	var creds []model.WebAuthnCredential
	result := db.conn.Where("tg_id = ?", tgID).Order("created_at DESC").Find(&creds)
	if result.Error != nil {
		return nil, result.Error
	}
	return creds, nil
}

// UpdateWebAuthnCredential updates an existing credential
func (db *DB) UpdateWebAuthnCredential(cred *model.WebAuthnCredential) error {
	return db.conn.Save(cred).Error
}

// UpdateWebAuthnCredentialWithSessionCleanup updates a credential and deletes the session in a transaction
func (db *DB) UpdateWebAuthnCredentialWithSessionCleanup(cred *model.WebAuthnCredential, sessionID string) error {
	return db.conn.Transaction(func(tx *gorm.DB) error {
		// Update credential
		if err := tx.Save(cred).Error; err != nil {
			return err
		}

		// Cleanup session (non-critical, don't fail transaction)
		_ = tx.Delete(&model.WebAuthnSession{}, "id = ?", sessionID).Error

		return nil
	})
}

// DeleteWebAuthnCredential deletes a credential by its ID
func (db *DB) DeleteWebAuthnCredential(credentialID []byte) error {
	return db.conn.Delete(&model.WebAuthnCredential{}, "id = ?", credentialID).Error
}

// WebAuthn Session Operations

// CreateWebAuthnSession creates a new WebAuthn session
func (db *DB) CreateWebAuthnSession(session *model.WebAuthnSession) error {
	return db.conn.Create(session).Error
}

// GetWebAuthnSession retrieves a session by its ID
func (db *DB) GetWebAuthnSession(sessionID string) (*model.WebAuthnSession, error) {
	var session model.WebAuthnSession
	result := db.conn.Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

// DeleteWebAuthnSession deletes a session by its ID
func (db *DB) DeleteWebAuthnSession(sessionID string) error {
	return db.conn.Delete(&model.WebAuthnSession{}, "id = ?", sessionID).Error
}

// CleanupExpiredWebAuthnSessions removes expired sessions
func (db *DB) CleanupExpiredWebAuthnSessions() error {
	now := time.Now().UTC()
	return db.conn.Where("expires_at < ?", now).Delete(&model.WebAuthnSession{}).Error
}

// GetUserWithWebAuthnCredentials retrieves a user with preloaded credentials
func (db *DB) GetUserWithWebAuthnCredentials(tgID int64) (*model.User, error) {
	var user model.User
	result := db.conn.Preload("Credentials").Where("tg_id = ?", tgID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
