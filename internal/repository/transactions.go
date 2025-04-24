package repository

import (
	"happypoor/internal/db"
	"happypoor/internal/model"
	"time"
)

type Transactions struct {
	DB *db.DB
}

func (r *Transactions) Add(transaction model.Transaction) error {
	return r.DB.CreateTransaction(&transaction)
}

func (r *Transactions) GetMonthlyTotalsCurrentYear(tgID int64) (map[int]map[model.TransactionType]float64, error) {
	year := time.Now().Year()
	return r.DB.GetMonthlyTotalsInYear(tgID, year)
}

func (r *Transactions) GetUserTransactionsByMonthPaginated(tgID int64, year, month, offset, limit int) ([]model.Transaction, int64, error) {
	return r.DB.GetUserTransactionsByMonthPaginated(tgID, year, month, offset, limit)
}

// GetMonthCategorizedTotals returns the transaction totals for each category for a specific month
func (r *Transactions) GetMonthCategorizedTotals(tgID int64, year int, month int) (map[model.TransactionType]map[model.TransactionCategory]float64, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month

	// Get expense categories
	expenseTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeExpense)
	if err != nil {
		return nil, err
	}

	// Get income categories
	incomeTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeIncome)
	if err != nil {
		return nil, err
	}

	// Return both in a map
	return map[model.TransactionType]map[model.TransactionCategory]float64{
		model.TypeExpense: expenseTotals,
		model.TypeIncome:  incomeTotals,
	}, nil
}

// GetYearCategorizedTotals returns the transaction totals for each category for a specific year
func (r *Transactions) GetYearCategorizedTotals(tgID int64, year int) (map[model.TransactionType]map[model.TransactionCategory]float64, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	// Get expense categories
	expenseTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeExpense)
	if err != nil {
		return nil, err
	}

	// Get income categories
	incomeTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeIncome)
	if err != nil {
		return nil, err
	}

	// Return both in a map
	return map[model.TransactionType]map[model.TransactionCategory]float64{
		model.TypeExpense: expenseTotals,
		model.TypeIncome:  incomeTotals,
	}, nil
}
