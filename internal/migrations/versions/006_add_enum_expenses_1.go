package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigration("006", "Add Pets enum to transactions", addEnumExpensesEnum1)
}

func addEnumExpensesEnum1(tx *gorm.DB) error {
	db, err := tx.DB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		ALTER TYPE transaction_category ADD VALUE IF NOT EXISTS 'Pets';
	`)

	return err
}
