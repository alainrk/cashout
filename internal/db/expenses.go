package db

import (
	"happypoor/internal/model"
	"time"
)

// CreateExpense creates a new expense record
func (db *DB) CreateExpense(expense *model.Expense) error {
	return db.conn.Create(expense).Error
}

// GetExpenseByID retrieves an expense by its ID
func (db *DB) GetExpenseByID(id int64) (*model.Expense, error) {
	var expense model.Expense
	result := db.conn.Where("id = ?", id).First(&expense)
	if result.Error != nil {
		return nil, result.Error
	}
	return &expense, nil
}

// UpdateExpense updates an existing expense
func (db *DB) UpdateExpense(expense *model.Expense) error {
	return db.conn.Save(expense).Error
}

// DeleteExpense deletes an expense by ID
func (db *DB) DeleteExpense(id int64) error {
	return db.conn.Delete(&model.Expense{}, id).Error
}

// GetUserExpenses retrieves all expenses for a user
func (db *DB) GetUserExpenses(tgID int64) ([]model.Expense, error) {
	var expenses []model.Expense
	result := db.conn.Where("tg_id = ?", tgID).Order("date DESC").Find(&expenses)
	if result.Error != nil {
		return nil, result.Error
	}
	return expenses, nil
}

// GetUserExpensesByDateRange retrieves expenses for a user within a date range
func (db *DB) GetUserExpensesByDateRange(tgID int64, startDate, endDate time.Time) ([]model.Expense, error) {
	var expenses []model.Expense
	result := db.conn.Where("tg_id = ? AND date BETWEEN ? AND ?",
		tgID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("date DESC").
		Find(&expenses)

	if result.Error != nil {
		return nil, result.Error
	}
	return expenses, nil
}

// GetUserExpensesByMonth retrieves expenses for a user for a specific year and month
func (db *DB) GetUserExpensesByMonth(tgID int64, year int, month int) ([]model.Expense, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month

	return db.GetUserExpensesByDateRange(tgID, startDate, endDate)
}

// GetUserExpensesByCategory retrieves expenses for a user grouped by category
func (db *DB) GetUserExpensesByCategory(tgID int64, startDate, endDate time.Time, expenseType model.ExpenseType) (map[model.ExpenseCategory]float64, error) {
	var results []struct {
		Category model.ExpenseCategory
		Total    float64
	}

	query := db.conn.Table("expenses").
		Select("category, SUM(amount) as total").
		Where("tg_id = ? AND date BETWEEN ? AND ? AND type = ?",
			tgID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), expenseType).
		Group("category").
		Order("total DESC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Convert to map
	categoryTotals := make(map[model.ExpenseCategory]float64)
	for _, result := range results {
		categoryTotals[result.Category] = result.Total
	}

	return categoryTotals, nil
}

// GetUserBalance calculates the total balance (income - expenses) for a user
func (db *DB) GetUserBalance(tgID int64, startDate, endDate time.Time) (float64, error) {
	var income float64
	var expense float64

	// Get total income
	incomeQuery := db.conn.Table("expenses").
		Select("COALESCE(SUM(amount), 0) as total").
		Where("tg_id = ? AND date BETWEEN ? AND ? AND type = ?",
			tgID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), model.TypeIncome)

	if err := incomeQuery.Scan(&income).Error; err != nil {
		return 0, err
	}

	// Get total expenses
	expenseQuery := db.conn.Table("expenses").
		Select("COALESCE(SUM(amount), 0) as total").
		Where("tg_id = ? AND date BETWEEN ? AND ? AND type = ?",
			tgID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), model.TypeExpense)

	if err := expenseQuery.Scan(&expense).Error; err != nil {
		return 0, err
	}

	return income - expense, nil
}

// GetMonthlyTotals gets monthly totals for a specific year
func (db *DB) GetMonthlyTotals(tgID int64, year int) (map[int]map[model.ExpenseType]float64, error) {
	var results []struct {
		Month int
		Type  model.ExpenseType
		Total float64
	}

	query := db.conn.Table("expenses").
		Select("EXTRACT(MONTH FROM date) as month, type, SUM(amount) as total").
		Where("tg_id = ? AND EXTRACT(YEAR FROM date) = ?", tgID, year).
		Group("month, type").
		Order("month")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Convert to map of maps: month -> type -> amount
	monthlyTotals := make(map[int]map[model.ExpenseType]float64)

	for _, result := range results {
		if _, exists := monthlyTotals[result.Month]; !exists {
			monthlyTotals[result.Month] = make(map[model.ExpenseType]float64)
		}
		monthlyTotals[result.Month][result.Type] = result.Total
	}

	return monthlyTotals, nil
}
