package versions

import (
	"cashout/internal/migrations"

	"gorm.io/gorm"
)

func init() {
	migrations.RegisterMigrationWithRollback("011", "Create budgets and budget_alerts tables", createBudgetsTables, rollbackBudgetsTables)
}

func createBudgetsTables(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS budgets (
			id          BIGSERIAL PRIMARY KEY,
			tg_id       BIGINT NOT NULL,
			amount      DECIMAL(15,2) NOT NULL CHECK (amount > 0),
			currency    currency_type NOT NULL DEFAULT 'EUR',
			created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_user_budget UNIQUE (tg_id)
		);

		CREATE INDEX IF NOT EXISTS idx_budgets_tg_id ON budgets (tg_id);

		ALTER TABLE budgets ADD CONSTRAINT fk_budgets_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id);

		CREATE TABLE IF NOT EXISTS budget_alerts (
			id          BIGSERIAL PRIMARY KEY,
			tg_id       BIGINT NOT NULL,
			year_month  CHAR(7) NOT NULL,
			threshold   SMALLINT NOT NULL CHECK (threshold IN (80, 100)),
			fired_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_user_month_threshold UNIQUE (tg_id, year_month, threshold)
		);

		CREATE INDEX IF NOT EXISTS idx_budget_alerts_tg_month ON budget_alerts (tg_id, year_month);

		ALTER TABLE budget_alerts ADD CONSTRAINT fk_budget_alerts_tg_id FOREIGN KEY (tg_id) REFERENCES users (tg_id);
	`).Error
}

func rollbackBudgetsTables(tx *gorm.DB) error {
	return tx.Exec(`
		DROP TABLE IF EXISTS budget_alerts;
		DROP TABLE IF EXISTS budgets;
	`).Error
}
