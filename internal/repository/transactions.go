package repository

import (
	"cashout/internal/model"
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Transactions struct {
	Repository
}

func (r *Transactions) Add(transaction model.Transaction) (err error) {
	// Assuming context.Background() for now as the function signature does not include context.
	// Ideally, context should be passed down from the caller.
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "add"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	err = r.DB.CreateTransaction(&transaction)
	return err
}

func (r *Transactions) GetByID(id int64) (transaction model.Transaction, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "get_by_id"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	dbTransaction, err := r.DB.GetTransactionByID(id)
	if err != nil {
		return model.Transaction{}, err
	}
	return *dbTransaction, nil
}

func (r *Transactions) Update(transaction *model.Transaction) (err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "update"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	err = r.DB.UpdateTransaction(transaction)
	return err
}

func (r *Transactions) Delete(id int64, tgID int64) (err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "delete"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	err = r.DB.DeleteTransactionByID(id, tgID)
	return err
}

func (r *Transactions) GetMonthlyTotalsCurrentYear(tgID int64) (result map[int]map[model.TransactionType]float64, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "get_monthly_totals_current_year"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	year := time.Now().Year()
	result, err = r.DB.GetMonthlyTotalsInYear(tgID, year)
	return result, err
}

func (r *Transactions) GetUserTransactionsByMonthPaginated(tgID int64, year, month, offset, limit int) (transactions []model.Transaction, total int64, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "get_user_transactions_by_month_paginated"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	transactions, total, err = r.DB.GetUserTransactionsByMonthPaginated(tgID, year, month, offset, limit)
	return transactions, total, err
}

func (r *Transactions) GetUserTransactionsByDateRangePaginated(tgID int64, startDate, endDate time.Time, offset, limit int) (transactions []model.Transaction, total int64, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "get_user_transactions_by_date_range_paginated"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	transactions, total, err = r.DB.GetUserTransactionsByDateRangePaginated(tgID, startDate, endDate, offset, limit)
	return transactions, total, err
}

// GetUserTransactionsPaginated retrieves all transactions for a user with pagination
func (r *Transactions) GetUserTransactionsPaginated(tgID int64, offset, limit int) (transactions []model.Transaction, total int64, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"),
			attribute.String("db.operation", "get_user_transactions_paginated"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	transactions, total, err = r.DB.GetUserTransactionsPaginated(tgID, offset, limit)
	return transactions, total, err
}

// GetMonthCategorizedTotals returns the transaction totals for each category for a specific month
func (r *Transactions) GetMonthCategorizedTotals(tgID int64, year int, month int) (result map[model.TransactionType]map[model.TransactionCategory]float64, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"), // This operation might involve multiple queries or complex logic
			attribute.String("db.operation", "get_month_categorized_totals"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of the month

	// Get expense categories
	expenseTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeExpense)
	if err != nil {
		return nil, err // err is captured by defer
	}

	// Get income categories
	incomeTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeIncome)
	if err != nil {
		return nil, err // err is captured by defer
	}

	// Return both in a map
	return map[model.TransactionType]map[model.TransactionCategory]float64{
		model.TypeExpense: expenseTotals,
		model.TypeIncome:  incomeTotals,
	}, nil
}

// GetYearCategorizedTotals returns the transaction totals for each category for a specific year
func (r *Transactions) GetYearCategorizedTotals(tgID int64, year int) (result map[model.TransactionType]map[model.TransactionCategory]float64, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "transactions"), // Similar to GetMonthCategorizedTotals, might be complex
			attribute.String("db.operation", "get_year_categorized_totals"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	// Get expense categories
	expenseTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeExpense)
	if err != nil {
		return nil, err // err is captured by defer
	}

	// Get income categories
	incomeTotals, err := r.DB.GetUserTransactionsByCategory(tgID, startDate, endDate, model.TypeIncome)
	if err != nil {
		return nil, err // err is captured by defer
	}

	// Return both in a map
	result = map[model.TransactionType]map[model.TransactionCategory]float64{
		model.TypeExpense: expenseTotals,
		model.TypeIncome:  incomeTotals,
	}
	return result, nil // err will be nil if both sub-queries succeed
}
