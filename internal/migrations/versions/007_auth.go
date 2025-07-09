package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigrationWithRollback("007", "Create auth tables", createAuthTables, rollbackAuthTables)
}

func createAuthTables(tx *gorm.DB) error {
	return tx.Exec(`
		-- Create auth status enum
		DROP TYPE IF EXISTS auth_status CASCADE;
		CREATE TYPE auth_status AS ENUM (
			'pending',
			'verified',
			'expired'
		);

		-- Create auth_tokens table
		CREATE TABLE IF NOT EXISTS auth_tokens (
			id SERIAL PRIMARY KEY,
			tg_id BIGINT NOT NULL,
			token VARCHAR(32) NOT NULL UNIQUE,
			status auth_status NOT NULL DEFAULT 'pending',
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		-- Create indexes for auth_tokens
		CREATE INDEX IF NOT EXISTS idx_auth_tokens_tg_id ON auth_tokens (tg_id);
		CREATE INDEX IF NOT EXISTS idx_auth_tokens_token ON auth_tokens (token);
		CREATE INDEX IF NOT EXISTS idx_auth_tokens_status ON auth_tokens (status);
		CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON auth_tokens (expires_at);
		
		-- Add foreign key constraint
		ALTER TABLE auth_tokens ADD CONSTRAINT fk_auth_tokens_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id) ON DELETE CASCADE;

		-- Create web_sessions table
		CREATE TABLE IF NOT EXISTS web_sessions (
			id VARCHAR(64) PRIMARY KEY,
			tg_id BIGINT NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		-- Create indexes for web_sessions
		CREATE INDEX IF NOT EXISTS idx_web_sessions_tg_id ON web_sessions (tg_id);
		CREATE INDEX IF NOT EXISTS idx_web_sessions_expires_at ON web_sessions (expires_at);
		
		-- Add foreign key constraint
		ALTER TABLE web_sessions ADD CONSTRAINT fk_web_sessions_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id) ON DELETE CASCADE;
	`).Error
}

func rollbackAuthTables(tx *gorm.DB) error {
	return tx.Exec(`
		DROP TABLE IF EXISTS web_sessions;
		DROP TABLE IF EXISTS auth_tokens;
		DROP TYPE IF EXISTS auth_status;
	`).Error
}
