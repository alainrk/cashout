package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigrationWithRollback("012", "Create api_tokens table", createAPITokensTable, rollbackAPITokensTable)
}

func createAPITokensTable(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS api_tokens (
			id           BIGSERIAL PRIMARY KEY,
			tg_id        BIGINT NOT NULL,
			name         VARCHAR(255) NOT NULL,
			token_hash   CHAR(64) NOT NULL UNIQUE,
			prefix       VARCHAR(16) NOT NULL,
			expires_at   TIMESTAMP WITH TIME ZONE NULL,
			last_used_at TIMESTAMP WITH TIME ZONE NULL,
			created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_api_tokens_tg_id ON api_tokens (tg_id);

		ALTER TABLE api_tokens ADD CONSTRAINT fk_api_tokens_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id);
	`).Error
}

func rollbackAPITokensTable(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS api_tokens;`).Error
}
