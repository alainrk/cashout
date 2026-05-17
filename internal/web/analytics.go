package web

import (
	"net/http"
	"strconv"
	"time"

	"cashout/internal/client"
	"cashout/internal/db"
	"cashout/internal/model"
)

// categoryEntry is the per-category breakdown returned by analytics endpoints.
type categoryEntry struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Count    int64   `json:"count"`
	Pct      float64 `json:"pct"`
}

// monthPoint is a single month's pivoted totals.
type monthPoint struct {
	Month   string  `json:"month"` // YYYY-MM
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

// handleAPIAnalyticsMonthly returns category breakdown + totals for a month.
func (s *Server) handleAPIAnalyticsMonthly(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	monthStr := r.URL.Query().Get("month")
	current, err := time.Parse(monthLayout, monthStr)
	if err != nil {
		current = time.Now()
	}

	startDate := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	expense, err := s.repositories.Transactions.GetCategoryAggregates(user.TgID, startDate, endDate, model.TypeExpense)
	if err != nil {
		s.logger.Errorf("analytics monthly expense: %v", err)
		s.sendJSONError(w, "Failed to load analytics", http.StatusInternalServerError)
		return
	}
	income, err := s.repositories.Transactions.GetCategoryAggregates(user.TgID, startDate, endDate, model.TypeIncome)
	if err != nil {
		s.logger.Errorf("analytics monthly income: %v", err)
		s.sendJSONError(w, "Failed to load analytics", http.StatusInternalServerError)
		return
	}

	expenseEntries, totalExpense := buildCategoryEntries(expense)
	incomeEntries, totalIncome := buildCategoryEntries(income)

	resp := map[string]any{
		"month":         current.Format(monthLayout),
		"totalIncome":   totalIncome,
		"totalExpenses": totalExpense,
		"balance":       totalIncome - totalExpense,
		"byCategory": map[string]any{
			"Expense": expenseEntries,
			"Income":  incomeEntries,
		},
	}
	s.sendJSONSuccess(w, resp)
}

// handleAPIAnalyticsTrend returns trailing N-month income/expense/balance.
func (s *Server) handleAPIAnalyticsTrend(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	months := 12
	if v := r.URL.Query().Get("months"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 60 {
			months = n
		}
	}

	now := time.Now().UTC()
	endMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
	startMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -(months - 1), 0)

	rows, err := s.repositories.Transactions.GetMonthlyTotalsByRange(user.TgID, startMonth, endMonth)
	if err != nil {
		s.logger.Errorf("analytics trend: %v", err)
		s.sendJSONError(w, "Failed to load trend", http.StatusInternalServerError)
		return
	}

	points := buildMonthPoints(startMonth, months, rows)

	s.sendJSONSuccess(w, map[string]any{
		"from":   points[0].Month,
		"to":     points[len(points)-1].Month,
		"points": points,
	})
}

// handleAPIAnalyticsYear returns per-month totals and category breakdown for a year.
func (s *Server) handleAPIAnalyticsYear(w http.ResponseWriter, r *http.Request) {
	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	year := time.Now().Year()
	if v := r.URL.Query().Get("year"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1970 && n <= 9999 {
			year = n
		}
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)

	rows, err := s.repositories.Transactions.GetMonthlyTotalsByRange(user.TgID, startDate, endDate)
	if err != nil {
		s.logger.Errorf("analytics year monthly: %v", err)
		s.sendJSONError(w, "Failed to load year analytics", http.StatusInternalServerError)
		return
	}
	points := buildMonthPoints(startDate, 12, rows)

	byMonth := make([]map[string]any, 12)
	var totalIncome, totalExpense float64
	for i, p := range points {
		byMonth[i] = map[string]any{
			"month":   i + 1,
			"income":  p.Income,
			"expense": p.Expense,
			"balance": p.Balance,
		}
		totalIncome += p.Income
		totalExpense += p.Expense
	}

	expense, err := s.repositories.Transactions.GetCategoryAggregates(user.TgID, startDate, endDate, model.TypeExpense)
	if err != nil {
		s.logger.Errorf("analytics year expense: %v", err)
		s.sendJSONError(w, "Failed to load year analytics", http.StatusInternalServerError)
		return
	}
	income, err := s.repositories.Transactions.GetCategoryAggregates(user.TgID, startDate, endDate, model.TypeIncome)
	if err != nil {
		s.logger.Errorf("analytics year income: %v", err)
		s.sendJSONError(w, "Failed to load year analytics", http.StatusInternalServerError)
		return
	}
	expenseEntries, _ := buildCategoryEntries(expense)
	incomeEntries, _ := buildCategoryEntries(income)

	s.sendJSONSuccess(w, map[string]any{
		"year":          year,
		"totalIncome":   totalIncome,
		"totalExpenses": totalExpense,
		"balance":       totalIncome - totalExpense,
		"byMonth":       byMonth,
		"byCategory": map[string]any{
			"Expense": expenseEntries,
			"Income":  incomeEntries,
		},
	})
}

func buildCategoryEntries(rows []db.CategoryAggregate) ([]categoryEntry, float64) {
	var total float64
	for _, r := range rows {
		total += r.Amount
	}
	entries := make([]categoryEntry, len(rows))
	for i, r := range rows {
		pct := 0.0
		if total > 0 {
			pct = (r.Amount / total) * 100
		}
		entries[i] = categoryEntry{
			Category: string(r.Category),
			Amount:   r.Amount,
			Count:    r.Count,
			Pct:      pct,
		}
	}
	return entries, total
}

func buildMonthPoints(start time.Time, months int, rows []db.MonthTotal) []monthPoint {
	idx := make(map[string]int, months)
	points := make([]monthPoint, months)
	for i := range months {
		ym := start.AddDate(0, i, 0).Format(monthLayout)
		points[i] = monthPoint{Month: ym}
		idx[ym] = i
	}
	for _, r := range rows {
		i, ok := idx[r.YM]
		if !ok {
			continue
		}
		switch r.Type {
		case model.TypeIncome:
			points[i].Income = r.Total
		case model.TypeExpense:
			points[i].Expense = r.Total
		}
	}
	for i := range points {
		points[i].Balance = points[i].Income - points[i].Expense
	}
	return points
}
