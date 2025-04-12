package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the database wrapper
type DB struct {
	conn *gorm.DB
}

// NewDB initializes a new database connection
func NewDB(postgresURL string) (*DB, error) {
	conn, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = conn.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return &DB{conn: conn}, nil
}

// GetUser retrieves a user by Telegram ID
func (db *DB) GetUser(tgID int64) (*User, error) {
	var user User
	result := db.conn.Where("tg_id = ?", tgID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by Telegram username
func (db *DB) GetUserByUsername(username string) (*User, error) {
	var user User
	result := db.conn.Where("tg_username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// SetUser creates or updates a user
func (db *DB) SetUser(user *User) error {
	// Use upsert functionality (create if not exists, update if exists)
	result := db.conn.Save(user)
	return result.Error
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
