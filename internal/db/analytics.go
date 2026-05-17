package db

import (
	"cashout/internal/model"
	"time"
)

// CategoryAggregate is the per-category aggregation used by analytics endpoints.
type CategoryAggregate struct {
	Category model.TransactionCategory
	Amount   float64
	Count    int64
}

// GetCategoryAggregates returns per-category totals and counts for a user/type
// over a date range (inclusive).
func (db *DB) GetCategoryAggregates(tgID int64, startDate, endDate time.Time, transactionType model.TransactionType) ([]CategoryAggregate, error) {
	var rows []struct {
		Category model.TransactionCategory
		Amount   float64
		Count    int64
	}

	err := db.conn.Table("transactions").
		Select("category, SUM(amount) as amount, COUNT(*) as count").
		Where("tg_id = ? AND date BETWEEN ? AND ? AND type = ?",
			tgID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), transactionType).
		Group("category").
		Order("amount DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]CategoryAggregate, len(rows))
	for i, r := range rows {
		result[i] = CategoryAggregate{Category: r.Category, Amount: r.Amount, Count: r.Count}
	}
	return result, nil
}

// MonthTotal is per-month totals keyed by YYYY-MM string.
type MonthTotal struct {
	YM    string
	Type  model.TransactionType
	Total float64
}

// GetMonthlyTotalsByRange returns per (month, type) totals between startDate
// (inclusive) and endDate (inclusive), ordered by month ascending.
func (db *DB) GetMonthlyTotalsByRange(tgID int64, startDate, endDate time.Time) ([]MonthTotal, error) {
	var rows []struct {
		YM    string
		Type  model.TransactionType
		Total float64
	}

	err := db.conn.Table("transactions").
		Select("to_char(date, 'YYYY-MM') as ym, type, SUM(amount) as total").
		Where("tg_id = ? AND date BETWEEN ? AND ?",
			tgID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("ym, type").
		Order("ym").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]MonthTotal, len(rows))
	for i, r := range rows {
		result[i] = MonthTotal{YM: r.YM, Type: r.Type, Total: r.Total}
	}
	return result, nil
}
