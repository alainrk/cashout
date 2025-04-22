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
