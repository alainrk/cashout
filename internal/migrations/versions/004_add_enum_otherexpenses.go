package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigration("004", "Add OtherExpenses enum to transactions", addOtherExpensesEnum)
}

func addOtherExpensesEnum(tx *gorm.DB) error {
	db, err := tx.DB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		ALTER TYPE transaction_category ADD VALUE IF NOT EXISTS 'OtherExpenses';
	`)

	return err
}
