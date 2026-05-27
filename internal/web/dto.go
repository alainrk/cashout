package web

import "time"

// ErrorResponse is the body returned by sendJSONError.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse is a generic success body with a human-readable message.
type MessageResponse struct {
	Message string `json:"message"`
}

// TransactionDTO is the JSON shape of a transaction returned by the API.
type TransactionDTO struct {
	ID          int64     `json:"id"`
	Date        time.Time `json:"date"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
}

// TransactionsResponse is the body of GET /api/transactions.
type TransactionsResponse struct {
	Transactions []TransactionDTO `json:"transactions"`
	Count        int              `json:"count"`
}

// CategoriesResponse is the body of GET /api/categories.
type CategoriesResponse struct {
	Categories []string `json:"categories"`
}

// CreateTransactionRequest is the body of POST /api/transactions/create.
type CreateTransactionRequest struct {
	Type        string  `json:"type"        example:"Expense"`
	Category    string  `json:"category"    example:"Food"`
	Amount      float64 `json:"amount"      example:"12.50"`
	Description string  `json:"description" example:"lunch"`
	Date        string  `json:"date"        example:"2026-05-21"`
}

// DeleteTransactionRequest is the body of DELETE /api/transactions/delete.
type DeleteTransactionRequest struct {
	ID int64 `json:"id"`
}

// EditTransactionRequest is the body of PATCH /api/transactions/edit.
// Only non-nil fields are applied. Type is intentionally not editable —
// switching between Income/Expense is forbidden, matching the Telegram bot.
type EditTransactionRequest struct {
	ID          int64    `json:"id"                    example:"42"`
	Category    *string  `json:"category,omitempty"    example:"Grocery"`
	Amount      *float64 `json:"amount,omitempty"      example:"19.90"`
	Description *string  `json:"description,omitempty" example:"weekly shop"`
	Date        *string  `json:"date,omitempty"        example:"2026-05-21"`
}

// CloneTransactionRequest is the body of POST /api/transactions/clone.
type CloneTransactionRequest struct {
	ID int64 `json:"id" example:"42"`
}

// SearchTransactionsRequest is the body of POST /api/transactions/search.
// Empty / zero fields are ignored. Type accepts "Income", "Expense" or "".
// Category "" or "all" disables the category filter.
type SearchTransactionsRequest struct {
	Query     string   `json:"query,omitempty"     example:"coffee"`
	Category  string   `json:"category,omitempty"  example:"Grocery"`
	Type      string   `json:"type,omitempty"      example:"Expense"`
	DateFrom  string   `json:"dateFrom,omitempty"  example:"2026-01-01"`
	DateTo    string   `json:"dateTo,omitempty"    example:"2026-05-31"`
	AmountMin *float64 `json:"amountMin,omitempty" example:"5"`
	AmountMax *float64 `json:"amountMax,omitempty" example:"100"`
	Offset    int      `json:"offset,omitempty"    example:"0"`
	Limit     int      `json:"limit,omitempty"     example:"50"`
}

// SearchTransactionsResponse is the body of POST /api/transactions/search.
type SearchTransactionsResponse struct {
	Transactions []TransactionDTO `json:"transactions"`
	Total        int64            `json:"total"`
	Offset       int              `json:"offset"`
	Limit        int              `json:"limit"`
}

// StatsResponse is the body of GET /api/stats.
type StatsResponse struct {
	Balance           float64 `json:"balance"`
	TotalIncome       float64 `json:"totalIncome"`
	TotalExpenses     float64 `json:"totalExpenses"`
	TotalTransactions int     `json:"totalTransactions"`
}

// BudgetResponse is the body of GET/POST/PUT/DELETE /api/budget.
// When HasBudget is false the other fields are zero values.
type BudgetResponse struct {
	HasBudget bool    `json:"hasBudget"`
	Amount    float64 `json:"amount,omitempty"`
	Currency  string  `json:"currency,omitempty"`
	Spent     float64 `json:"spent,omitempty"`
	Pct       int     `json:"pct,omitempty"`
	Month     string  `json:"month,omitempty"`
}

// BudgetUpsertRequest is the body of POST/PUT /api/budget.
type BudgetUpsertRequest struct {
	Amount float64 `json:"amount"`
}

// CategoryEntry is one row of a category breakdown.
type CategoryEntry struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Count    int64   `json:"count"`
	Pct      float64 `json:"pct"`
}

// CategoryBreakdown groups category entries by transaction type.
type CategoryBreakdown struct {
	Expense []CategoryEntry `json:"Expense"`
	Income  []CategoryEntry `json:"Income"`
}

// MonthPoint is one month's pivoted totals.
type MonthPoint struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

// MonthlyAnalyticsResponse is the body of GET /api/analytics/monthly.
type MonthlyAnalyticsResponse struct {
	Month         string            `json:"month"`
	TotalIncome   float64           `json:"totalIncome"`
	TotalExpenses float64           `json:"totalExpenses"`
	Balance       float64           `json:"balance"`
	ByCategory    CategoryBreakdown `json:"byCategory"`
}

// TrendResponse is the body of GET /api/analytics/trend.
type TrendResponse struct {
	From   string       `json:"from"`
	To     string       `json:"to"`
	Points []MonthPoint `json:"points"`
}

// YearMonthEntry is one month inside the year analytics breakdown.
type YearMonthEntry struct {
	Month   int     `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

// YearAnalyticsResponse is the body of GET /api/analytics/year.
type YearAnalyticsResponse struct {
	Year          int               `json:"year"`
	TotalIncome   float64           `json:"totalIncome"`
	TotalExpenses float64           `json:"totalExpenses"`
	Balance       float64           `json:"balance"`
	ByMonth       []YearMonthEntry  `json:"byMonth"`
	ByCategory    CategoryBreakdown `json:"byCategory"`
}
