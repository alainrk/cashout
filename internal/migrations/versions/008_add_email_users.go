package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigrationWithRollback("008", "Add email to users", addEmailUsers, rollbackEmailUsers)
}

func addEmailUsers(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT UNIQUE NULL;
	`).Error
}

func rollbackEmailUsers(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE users DROP COLUMN email IF EXISTS;
	`).Error
}
