package db

import (
	"fmt"
	"happypoor/internal/model"

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
	err = conn.AutoMigrate(&model.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
