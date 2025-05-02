package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigration("002", "Add name to users", addNameUsers)
}

func addNameUsers(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE users ADD COLUMN name TEXT;
	`).Error
}
