package db

import (
	"time"

	"cashout/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpsertBudget inserts or updates the single budget row for a user.
func (db *DB) UpsertBudget(budget *model.Budget) error {
	return db.conn.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "tg_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"amount", "currency", "updated_at",
		}),
	}).Create(budget).Error
}

// DeleteBudget removes the budget for a user. Returns gorm.ErrRecordNotFound if none.
func (db *DB) DeleteBudget(tgID int64) error {
	result := db.conn.Where("tg_id = ?", tgID).Delete(&model.Budget{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetBudget returns the user's budget or gorm.ErrRecordNotFound.
func (db *DB) GetBudget(tgID int64) (*model.Budget, error) {
	var b model.Budget
	err := db.conn.Where("tg_id = ?", tgID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// GetTotalExpensesForMonth sums all Expense amounts for a user in a calendar month.
func (db *DB) GetTotalExpensesForMonth(tgID int64, year, month int) (float64, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var total float64
	err := db.conn.Table("transactions").
		Select("COALESCE(SUM(amount), 0) as total").
		Where("tg_id = ? AND date BETWEEN ? AND ? AND type = ?",
			tgID,
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
			model.TypeExpense,
		).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// TryMarkAlertFired inserts an alert row; returns true if the insert actually happened
// (i.e. the alert had not yet fired for this user/month/threshold).
func (db *DB) TryMarkAlertFired(tgID int64, yearMonth string, threshold int16) (bool, error) {
	alert := model.BudgetAlert{
		TgID:      tgID,
		YearMonth: yearMonth,
		Threshold: threshold,
	}

	result := db.conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tg_id"}, {Name: "year_month"}, {Name: "threshold"}},
		DoNothing: true,
	}).Create(&alert)

	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}
