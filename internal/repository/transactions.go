package repository

import (
	"happypoor/internal/db"
	"happypoor/internal/model"
)

type Transactions struct {
	DB *db.DB
}

func (r *Transactions) Add(transaction model.Transaction) error {
	return r.DB.CreateTransaction(&transaction)
}
