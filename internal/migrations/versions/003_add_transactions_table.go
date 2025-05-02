package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigration("003", "Add transactions table", addTransactionsTable)
}

func addTransactionsTable(tx *gorm.DB) error {
	return tx.Exec(`
		-- Create transaction category enum type
		DROP TYPE IF EXISTS transaction_category;
		CREATE TYPE transaction_category AS ENUM (
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

		-- Create transaction type enum
		DROP TYPE IF EXISTS transaction_type;
		CREATE TYPE transaction_type AS ENUM (
			'Income',
			'Expense'
		);

		-- Create currency enum
		DROP TYPE IF EXISTS currency_type;
		CREATE TYPE currency_type AS ENUM (
			'EUR',
			'USD',
			'GBP',
			'JPY',
			'CHF'
		);

		-- Create transactions table
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			tg_id BIGINT NOT NULL,
			date DATE NOT NULL DEFAULT CURRENT_DATE,
			type transaction_type NOT NULL,
			category transaction_category NOT NULL,
			amount DECIMAL(15, 2) NOT NULL,
			currency currency_type NOT NULL DEFAULT 'EUR',
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		-- Create indexes
		CREATE INDEX IF NOT EXISTS idx_transactions_tg_id ON transactions (tg_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions (date);
		CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions (type);
		CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions (category);
		
		-- Add foreign key constraint
		ALTER TABLE transactions ADD CONSTRAINT fk_transactions_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id);
	`).Error
}

// Rollback function if needed
func rollbacktransactionsTable(tx *gorm.DB) error {
	return tx.Exec(`
		DROP TABLE IF EXISTS transactions;
		DROP TYPE IF EXISTS transaction_category;
		DROP TYPE IF EXISTS transaction_type;
		DROP TYPE IF EXISTS currency_type;
	`).Error
}
