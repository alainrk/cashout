// Package db provides a wrapper around the database
package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database wrapper
type DB struct {
	conn *gorm.DB
}

// NewDB initializes a new database connection
func NewDB(postgresURL string) (*DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	conn, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{
		// Disable foreign key constraints during AutoMigrate
		// We handle foreign keys explicitly in our migration files
		DisableForeignKeyConstraintWhenMigrating: true,
		// ErrRecordNotFound is used as control flow throughout the codebase
		// (e.g. "does the user have a budget yet?"); don't log it as a warning.
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Transaction wraps a database transaction (exposes GORM's transaction method)
func (db *DB) Transaction(fn func(tx *gorm.DB) error) error {
	return db.conn.Transaction(fn)
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping checks the database connection
func (db *DB) Ping() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
