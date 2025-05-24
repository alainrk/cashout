package repository

import (
	"cashout/internal/model"
	"time"
)

type Transactions struct {
	Repository
}

func (r *Transactions) Add(transaction model.Transaction) error {
	return r.DB.CreateTransaction(&transaction)
}

func (r *Transactions) GetByID(id int64) (model.Transaction, error) {
	transaction, err := r.DB.GetTransactionByID(id)
	if err != nil {
		return model.Transaction{}, err
	}
	return *transaction, nil
}

func (r *Transactions) Update(transaction *model.Transaction) error {
	return r.DB.UpdateTransaction(transaction)
}

func (r *Transactions) Delete(id int64, tgID int64) error {
	return r.DB.DeleteTransactionByID(id, tgID)
}

func (r *Transactions) GetMonthlyTotalsCurrentYear(tgID int64) (map[int]map[model.TransactionType]float64, error) {
	year := time.Now().Year()
	return r.DB.GetMonthlyTotalsInYear(tgID, year)
}

func (r *Transactions) GetUserTransactionsByMonthPaginated(tgID int64, year, month, offset, limit int) ([]model.Transaction, int64, error) {
	return r.DB.GetUserTransactionsByMonthPaginated(tgID, year, month, offset, limit)
}

func (r *Transactions) GetUserTransactionsByDateRangePaginated(tgID int64, startDate, endDate time.Time, offset, limit int) ([]model.Transaction, int64, error) {
	return r.DB.GetUserTransactionsByDateRangePaginated(tgID, startDate, endDate, offset, limit)
}

// GetUserTransactionsPaginated retrieves all transactions for a user with pagination
func (r *Transactions) GetUserTransactionsPaginated(tgID int64, offset, limit int) ([]model.Transaction, int64, error) {
	return r.DB.GetUserTransactionsPaginated(tgID, offset, limit)
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

// GetUserTransactions retrieves all transactions for a user (no pagination)
func (r *Transactions) GetUserTransactions(tgID int64) ([]model.Transaction, error) {
	return r.DB.GetUserTransactions(tgID)
}

// GetUserTransactionsByDateRange retrieves transactions for a user within a date range
func (r *Transactions) GetUserTransactionsByDateRange(tgID int64, startDate, endDate time.Time) ([]model.Transaction, error) {
	return r.DB.GetUserTransactionsByDateRange(tgID, startDate, endDate)
}
