package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigrationWithRollback("009", "Add partial unique index on email", addUniqueEmailIndex, rollbackUniqueEmailIndex)
}

func addUniqueEmailIndex(tx *gorm.DB) error {
	// Drop old constraint if it exists
	_ = tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key`).Error
	_ = tx.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_email`).Error

	// Create partial unique index (only enforces uniqueness on non-null emails)
	return tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
		ON users(email)
		WHERE email IS NOT NULL
	`).Error
}

func rollbackUniqueEmailIndex(tx *gorm.DB) error {
	return tx.Exec(`
		DROP INDEX IF EXISTS idx_users_email
	`).Error
}
