package web

import (
	"html/template"
	"net/http"
	"time"

	"cashout/internal/client"
	"cashout/internal/model"
)

const (
	monthLayout = "2006-01"
)

// handleDashboard shows the main dashboard
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, basePath+"/login", http.StatusSeeOther)
		return
	}

	// Parse month from query, default to current month
	monthStr := r.URL.Query().Get("month")
	currentMonth, err := time.Parse(monthLayout, monthStr)
	if err != nil || currentMonth.After(time.Now()) {
		currentMonth = time.Now()
	}

	// Calculate previous and next months
	prevMonth := currentMonth.AddDate(0, -1, 0)
	nextMonth := currentMonth.AddDate(0, 1, 0)

	// Disable next month button if it's the future
	now := time.Now()
	isCurrentMonth := currentMonth.Format(monthLayout) == now.Format(monthLayout)

	t, err := template.ParseFiles("web/templates/dashboard.html")
	if err != nil {
		s.logger.Errorf("Failed to parse template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User              *model.User
		CurrentMonthTitle string
		CurrentMonth      string
		PrevMonth         string
		NextMonth         string
		IsCurrentMonth    bool
	}{
		User:              user,
		CurrentMonthTitle: currentMonth.Format("January 2006"),
		CurrentMonth:      currentMonth.Format(monthLayout),
		PrevMonth:         prevMonth.Format(monthLayout),
		NextMonth:         nextMonth.Format(monthLayout),
		IsCurrentMonth:    isCurrentMonth,
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		s.logger.Errorf("Failed to execute template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

// handleAPIStats returns user statistics for a given month
func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse month from query, default to current month
	monthStr := r.URL.Query().Get("month")
	currentMonth, err := time.Parse(monthLayout, monthStr)
	if err != nil {
		currentMonth = time.Now()
	}

	// Get transactions for the month
	startDate := time.Date(currentMonth.Year(), currentMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)
	transactions, err := s.repositories.Transactions.GetUserTransactionsByDateRange(user.TgID, startDate, endDate)
	if err != nil {
		s.sendJSONError(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Calculate statistics
	var totalIncome, totalExpenses float64

	for _, tx := range transactions {
		if tx.Type == model.TypeIncome {
			totalIncome += tx.Amount
		} else {
			totalExpenses += tx.Amount
		}
	}

	balance := totalIncome - totalExpenses

	stats := map[string]any{
		"balance":           balance,
		"totalIncome":       totalIncome,
		"totalExpenses":     totalExpenses,
		"totalTransactions": len(transactions),
	}

	s.sendJSONSuccess(w, stats)
}

// handleAPITransactions returns user transactions for a given month
func (s *Server) handleAPITransactions(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse month from query, default to current month
	monthStr := r.URL.Query().Get("month")
	currentMonth, err := time.Parse(monthLayout, monthStr)
	if err != nil {
		currentMonth = time.Now()
	}

	// Get transactions for the month
	startDate := time.Date(currentMonth.Year(), currentMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)
	transactions, err := s.repositories.Transactions.GetUserTransactionsByDateRange(user.TgID, startDate, endDate)
	if err != nil {
		s.sendJSONError(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	type TransactionResponse struct {
		ID          int64     `json:"id"`
		Date        time.Time `json:"date"`
		Category    string    `json:"category"`
		Description string    `json:"description"`
		Amount      float64   `json:"amount"`
		Type        string    `json:"type"`
	}

	transactionResponses := make([]TransactionResponse, len(transactions))
	for i, tx := range transactions {
		transactionResponses[i] = TransactionResponse{
			ID:          tx.ID,
			Date:        tx.Date,
			Category:    string(tx.Category),
			Description: tx.Description,
			Amount:      tx.Amount,
			Type:        string(tx.Type),
		}
	}

	response := map[string]any{
		"transactions": transactionResponses,
		"count":        len(transactionResponses),
	}

	s.sendJSONSuccess(w, response)
}
