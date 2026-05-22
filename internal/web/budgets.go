package web

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"time"

	"cashout/internal/client"
	"cashout/internal/model"

	"gorm.io/gorm"
)

// handleAPIBudget multiplexes GET/POST/PUT/DELETE on /api/budget.
//
//	@Summary		Get the current monthly budget
//	@Tags			budget
//	@Produce		json
//	@Success		200	{object}	BudgetResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/budget [get]
//
//	@Summary		Create or update the monthly budget
//	@Tags			budget
//	@Accept			json
//	@Produce		json
//	@Param			body	body		BudgetUpsertRequest	true	"Budget amount"
//	@Success		200		{object}	BudgetResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/budget [post]
//
//	@Summary		Replace the monthly budget
//	@Tags			budget
//	@Accept			json
//	@Produce		json
//	@Param			body	body		BudgetUpsertRequest	true	"Budget amount"
//	@Success		200		{object}	BudgetResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/budget [put]
//
//	@Summary		Delete the monthly budget
//	@Tags			budget
//	@Produce		json
//	@Success		200	{object}	BudgetResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/budget [delete]
func (s *Server) handleAPIBudget(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.budgetGet(w, user)
	case http.MethodPost, http.MethodPut:
		s.budgetUpsert(w, r, user)
	case http.MethodDelete:
		s.budgetDelete(w, user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) budgetGet(w http.ResponseWriter, user *model.User) {
	budget, err := s.repositories.Budgets.Get(user.TgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.sendJSONSuccess(w, BudgetResponse{HasBudget: false})
			return
		}
		s.logger.Errorf("Failed to get budget: %v", err)
		s.sendJSONError(w, "Failed to get budget", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	spent, err := s.repositories.Budgets.TotalExpensesForMonth(user.TgID, now.Year(), int(now.Month()))
	if err != nil {
		s.logger.Errorf("Failed to compute month total: %v", err)
		s.sendJSONError(w, "Failed to compute progress", http.StatusInternalServerError)
		return
	}
	pct := 0
	if budget.Amount > 0 {
		pct = int(math.Floor(spent / budget.Amount * 100))
	}

	s.sendJSONSuccess(w, BudgetResponse{
		HasBudget: true,
		Amount:    budget.Amount,
		Currency:  string(budget.Currency),
		Spent:     spent,
		Pct:       pct,
		Month:     now.Format("2006-01"),
	})
}

func (s *Server) budgetUpsert(w http.ResponseWriter, r *http.Request, user *model.User) {
	var req BudgetUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		s.sendJSONError(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	budget := model.Budget{
		TgID:     user.TgID,
		Amount:   req.Amount,
		Currency: model.CurrencyEUR,
	}
	if err := s.repositories.Budgets.Upsert(&budget); err != nil {
		s.logger.Errorf("Failed to upsert budget: %v", err)
		s.sendJSONError(w, "Failed to save budget", http.StatusInternalServerError)
		return
	}

	s.budgetGet(w, user)
}

func (s *Server) budgetDelete(w http.ResponseWriter, user *model.User) {
	if err := s.repositories.Budgets.Delete(user.TgID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.sendJSONSuccess(w, BudgetResponse{HasBudget: false})
			return
		}
		s.logger.Errorf("Failed to delete budget: %v", err)
		s.sendJSONError(w, "Failed to delete budget", http.StatusInternalServerError)
		return
	}
	s.sendJSONSuccess(w, BudgetResponse{HasBudget: false})
}
