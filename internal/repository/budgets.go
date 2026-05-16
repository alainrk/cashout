package repository

import "cashout/internal/model"

type Budgets struct {
	Repository
}

func (r *Budgets) Upsert(budget *model.Budget) error {
	return r.DB.UpsertBudget(budget)
}

func (r *Budgets) Delete(tgID int64) error {
	return r.DB.DeleteBudget(tgID)
}

func (r *Budgets) Get(tgID int64) (*model.Budget, error) {
	return r.DB.GetBudget(tgID)
}

func (r *Budgets) TotalExpensesForMonth(tgID int64, year, month int) (float64, error) {
	return r.DB.GetTotalExpensesForMonth(tgID, year, month)
}

func (r *Budgets) TryMarkAlertFired(tgID int64, yearMonth string, threshold int16) (bool, error) {
	return r.DB.TryMarkAlertFired(tgID, yearMonth, threshold)
}
