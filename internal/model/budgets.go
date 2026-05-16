package model

import "time"

// Budget represents a single user's total monthly expense limit.
type Budget struct {
	ID        int64        `gorm:"column:id;primaryKey;autoIncrement"`
	TgID      int64        `gorm:"column:tg_id;not null;uniqueIndex"`
	Amount    float64      `gorm:"column:amount;not null;type:decimal(15,2)"`
	Currency  CurrencyType `gorm:"column:currency;not null;type:currency_type;default:'EUR'"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time    `gorm:"column:updated_at;autoUpdateTime"`
}

func (Budget) TableName() string {
	return "budgets"
}

// BudgetAlert tracks one-shot alert firings per (user, month, threshold).
type BudgetAlert struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	TgID      int64     `gorm:"column:tg_id;not null;index"`
	YearMonth string    `gorm:"column:year_month;not null;size:7"`
	Threshold int16     `gorm:"column:threshold;not null"`
	FiredAt   time.Time `gorm:"column:fired_at;autoCreateTime"`
}

func (BudgetAlert) TableName() string {
	return "budget_alerts"
}
