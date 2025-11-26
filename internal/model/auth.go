package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

// AuthStatus represents the status of an auth token
type AuthStatus string

const (
	AuthStatusPending  AuthStatus = "pending"
	AuthStatusVerified AuthStatus = "verified"
	AuthStatusExpired  AuthStatus = "expired"
)

// Value implements the driver.Valuer interface for AuthStatus
func (s AuthStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the sql.Scanner interface for AuthStatus
func (s *AuthStatus) Scan(value any) error {
	if value == nil {
		return errors.New("auth status cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid auth status")
	}

	*s = AuthStatus(strVal)
	return nil
}

// AuthToken represents the auth_tokens table structure
type AuthToken struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement"`
	TgID      int64      `gorm:"column:tg_id;not null;index"`
	Token     string     `gorm:"column:token;not null;unique;index"`
	Status    AuthStatus `gorm:"column:status;not null;type:auth_status;default:'pending';index"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null;index"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`

	// Association to User (optional)
	User *User `gorm:"foreignKey:TgID;references:TgID"`
}

// TableName overrides the table name
func (AuthToken) TableName() string {
	return "auth_tokens"
}

// IsValid checks if the token is still valid
func (a *AuthToken) IsValid() bool {
	return a.Status == AuthStatusPending && time.Now().UTC().Before(a.ExpiresAt)
}

// WebSession represents a web session
type WebSession struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TgID      int64     `gorm:"column:tg_id;not null;index"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null;index"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Association to User (optional)
	User *User `gorm:"foreignKey:TgID;references:TgID"`
}

// TableName overrides the table name
func (WebSession) TableName() string {
	return "web_sessions"
}

// IsValid checks if the session is still valid
func (s *WebSession) IsValid() bool {
	return time.Now().UTC().Before(s.ExpiresAt)
}
