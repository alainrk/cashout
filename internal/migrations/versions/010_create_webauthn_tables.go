package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigrationWithRollback("010", "Create WebAuthn tables", createWebAuthnTables, rollbackWebAuthnTables)
}

func createWebAuthnTables(tx *gorm.DB) error {
	return tx.Exec(`
		-- Create webauthn_credentials table
		CREATE TABLE IF NOT EXISTS webauthn_credentials (
			id BYTEA PRIMARY KEY,
			tg_id BIGINT NOT NULL,
			public_key BYTEA NOT NULL,
			attestation_type TEXT NOT NULL,
			aaguid BYTEA NOT NULL,
			sign_count INTEGER NOT NULL DEFAULT 0,
			clone_warning BOOLEAN NOT NULL DEFAULT false,
			flags_user_present BOOLEAN NOT NULL,
			flags_user_verified BOOLEAN NOT NULL,
			flags_backup_eligible BOOLEAN NOT NULL,
			flags_backup_state BOOLEAN NOT NULL,
			transport TEXT[] NOT NULL DEFAULT '{}',
			credential_name TEXT,
			last_used_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		-- Create indexes for webauthn_credentials
		CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_tg_id ON webauthn_credentials (tg_id);
		CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_last_used ON webauthn_credentials (last_used_at);

		-- Add foreign key constraint
		ALTER TABLE webauthn_credentials ADD CONSTRAINT fk_webauthn_credentials_tg_id
			FOREIGN KEY (tg_id) REFERENCES users (tg_id) ON DELETE CASCADE;

		-- Create webauthn_sessions table
		CREATE TABLE IF NOT EXISTS webauthn_sessions (
			id VARCHAR(64) PRIMARY KEY,
			tg_id BIGINT NOT NULL,
			challenge TEXT NOT NULL,
			user_verification TEXT NOT NULL,
			ceremony_type TEXT NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		-- Create indexes for webauthn_sessions
		CREATE INDEX IF NOT EXISTS idx_webauthn_sessions_tg_id ON webauthn_sessions (tg_id);
		CREATE INDEX IF NOT EXISTS idx_webauthn_sessions_expires_at ON webauthn_sessions (expires_at);

		-- Add foreign key constraint
		ALTER TABLE webauthn_sessions ADD CONSTRAINT fk_webauthn_sessions_tg_id
			FOREIGN KEY (tg_id) REFERENCES users (tg_id) ON DELETE CASCADE;
	`).Error
}

func rollbackWebAuthnTables(tx *gorm.DB) error {
	return tx.Exec(`
		DROP TABLE IF EXISTS webauthn_sessions;
		DROP TABLE IF EXISTS webauthn_credentials;
	`).Error
}
