package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigrationWithRollback("005", "Create reminders table", createRemindersTable, rollbackRemindersTable)
}

func createRemindersTable(tx *gorm.DB) error {
	return tx.Exec(`
		-- Create reminder type enum
		DROP TYPE IF EXISTS reminder_type;
		CREATE TYPE reminder_type AS ENUM (
			'weekly_recap',
			'monthly_recap',
			'yearly_recap'
		);

		-- Create reminder status enum
		DROP TYPE IF EXISTS reminder_status;
		CREATE TYPE reminder_status AS ENUM (
			'pending',
			'processing',
			'sent',
			'failed'
		);

		-- Create reminders table
		CREATE TABLE IF NOT EXISTS reminders (
			id SERIAL PRIMARY KEY,
			tg_id BIGINT NOT NULL,
			type reminder_type NOT NULL,
			status reminder_status NOT NULL DEFAULT 'pending',
			scheduled_for TIMESTAMP WITH TIME ZONE NOT NULL,
			processed_at TIMESTAMP WITH TIME ZONE,
			error_message TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			
			-- Ensure we don't create duplicate reminders for the same user/type/scheduled time
			CONSTRAINT unique_user_type_schedule UNIQUE (tg_id, type, scheduled_for)
		);

		-- Create indexes
		CREATE INDEX IF NOT EXISTS idx_reminders_tg_id ON reminders (tg_id);
		CREATE INDEX IF NOT EXISTS idx_reminders_status ON reminders (status);
		CREATE INDEX IF NOT EXISTS idx_reminders_scheduled_for ON reminders (scheduled_for);
		CREATE INDEX IF NOT EXISTS idx_reminders_type ON reminders (type);
		
		-- Add foreign key constraint
		ALTER TABLE reminders ADD CONSTRAINT fk_reminders_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id);
	`).Error
}

func rollbackRemindersTable(tx *gorm.DB) error {
	return tx.Exec(`
		DROP TABLE IF EXISTS reminders;
		DROP TYPE IF EXISTS reminder_type;
		DROP TYPE IF EXISTS reminder_status;
	`).Error
}
