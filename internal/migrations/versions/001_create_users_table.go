package versions

import (
	"happypoor/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigration("001", "Create users table", createUsersTable)
}

func createUsersTable(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			tg_id BIGINT PRIMARY KEY,
			tg_username VARCHAR(255) UNIQUE,
			tg_firstname VARCHAR(255),
			tg_lastname VARCHAR(255),
			session JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}
