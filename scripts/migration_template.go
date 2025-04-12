package versions

import (
	"myproject/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigration("{{VERSION}}", "{{NAME}}", {{FUNC_NAME}})
}

func {{FUNC_NAME}}(tx *gorm.DB) error {
	return tx.Exec(`
		-- Your SQL here
		-- For example:
		-- ALTER TABLE users ADD COLUMN new_column TEXT;
	`).Error
}
