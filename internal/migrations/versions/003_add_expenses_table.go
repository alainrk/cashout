package versions

import (
	"happypoor/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigration("003", "Add expenses table", addExpensesTable)
}

func addExpensesTable(tx *gorm.DB) error {
	return tx.Exec(`
		-- Create expense category enum type
		CREATE TYPE expense_category AS ENUM (
			'Salary',
			'OtherIncomes',
			'Car', 
			'Clothes', 
			'Grocery', 
			'House', 
			'Bills',
			'Entertainment',
			'Sport',
			'EatingOut', 
			'Transport', 
			'Learning',
			'Toiletry', 
			'Health', 
			'Tech', 
			'Gifts', 
			'Travel'
		);

		-- Create expense type enum
		CREATE TYPE expense_type AS ENUM (
			'Income',
			'Expense'
		);

		-- Create currency enum
		CREATE TYPE currency_type AS ENUM (
			'EUR',
			'USD',
			'GBP',
			'JPY',
			'CHF'
		);

		-- Create expenses table
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			tg_id BIGINT NOT NULL,
			date DATE NOT NULL DEFAULT CURRENT_DATE,
			type expense_type NOT NULL,
			category expense_category NOT NULL,
			amount DECIMAL(15, 2) NOT NULL,
			currency currency_type NOT NULL DEFAULT 'EUR',
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		-- Create indexes
		CREATE INDEX IF NOT EXISTS idx_expenses_tg_id ON expenses (tg_id);
		CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses (date);
		CREATE INDEX IF NOT EXISTS idx_expenses_type ON expenses (type);
		CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses (category);
		
		-- Add foreign key constraint
		ALTER TABLE expenses ADD CONSTRAINT fk_expenses_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id);
	`).Error
}

// Rollback function if needed
func rollbackExpensesTable(tx *gorm.DB) error {
	return tx.Exec(`
		DROP TABLE IF EXISTS expenses;
		DROP TYPE IF EXISTS expense_category;
		DROP TYPE IF EXISTS expense_type;
		DROP TYPE IF EXISTS currency_type;
	`).Error
}
