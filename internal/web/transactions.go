package web

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cashout/internal/client"
	"cashout/internal/model"
	"cashout/internal/repository"
)

const dateLayout = "2006-01-02"

func toTransactionDTO(tx model.Transaction) TransactionDTO {
	return TransactionDTO{
		ID:          tx.ID,
		Date:        tx.Date,
		Category:    string(tx.Category),
		Description: tx.Description,
		Amount:      tx.Amount,
		Type:        string(tx.Type),
	}
}

func isIncomeCategory(cat string) bool {
	for _, c := range model.GetIncomeCategories() {
		if c == cat {
			return true
		}
	}
	return false
}

// handleAPIEditTransaction applies a partial update to a transaction.
//
//	@Summary		Edit transaction (partial)
//	@Description	Update one or more fields of an existing transaction. Type cannot be changed; category must remain within the same type (Income↔Expense swaps are rejected).
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		EditTransactionRequest	true	"Fields to update; only non-null fields are applied"
//	@Success		200		{object}	TransactionDTO
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/transactions/edit [patch]
func (s *Server) handleAPIEditTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req EditTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.ID <= 0 {
		s.sendJSONError(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	tx, err := s.repositories.Transactions.GetByID(req.ID)
	if err != nil {
		s.sendJSONError(w, "Transaction not found", http.StatusNotFound)
		return
	}
	if tx.TgID != user.TgID {
		s.sendJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	if req.Category == nil && req.Amount == nil && req.Description == nil && req.Date == nil {
		s.sendJSONError(w, "No fields to update", http.StatusBadRequest)
		return
	}

	if req.Category != nil {
		cat := *req.Category
		if !model.IsValidTransactionCategory(cat) {
			s.sendJSONError(w, "Invalid category", http.StatusBadRequest)
			return
		}
		txIsIncome := tx.Type == model.TypeIncome
		if txIsIncome != isIncomeCategory(cat) {
			s.sendJSONError(w, "Cannot change between income and expense categories", http.StatusBadRequest)
			return
		}
		tx.Category = model.TransactionCategory(cat)
	}

	if req.Amount != nil {
		if *req.Amount <= 0 {
			s.sendJSONError(w, "Amount must be greater than 0", http.StatusBadRequest)
			return
		}
		tx.Amount = *req.Amount
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		if desc == "" {
			s.sendJSONError(w, "Description cannot be empty", http.StatusBadRequest)
			return
		}
		tx.Description = desc
	}

	if req.Date != nil {
		d, err := time.Parse(dateLayout, *req.Date)
		if err != nil {
			s.sendJSONError(w, "Invalid date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		if d.After(time.Now()) {
			s.sendJSONError(w, "Date cannot be in the future", http.StatusBadRequest)
			return
		}
		tx.Date = d
	}

	if err := s.repositories.Transactions.Update(&tx); err != nil {
		s.logger.Errorf("Failed to update transaction: %v", err)
		s.sendJSONError(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	s.sendJSONSuccess(w, toTransactionDTO(tx))
}

// handleAPICloneTransaction duplicates an existing transaction with today's date.
//
//	@Summary		Clone transaction
//	@Description	Duplicate an existing transaction; the new transaction copies type, category, amount, description and currency, with date set to today.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CloneTransactionRequest	true	"ID of the transaction to clone"
//	@Success		200		{object}	TransactionDTO
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/transactions/clone [post]
func (s *Server) handleAPICloneTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CloneTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.ID <= 0 {
		s.sendJSONError(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	source, err := s.repositories.Transactions.GetByID(req.ID)
	if err != nil {
		s.sendJSONError(w, "Transaction not found", http.StatusNotFound)
		return
	}
	if source.TgID != user.TgID {
		s.sendJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	clone := model.Transaction{
		TgID:        user.TgID,
		Date:        time.Now(),
		Type:        source.Type,
		Category:    source.Category,
		Amount:      source.Amount,
		Currency:    source.Currency,
		Description: source.Description,
	}

	if err := s.repositories.Transactions.Add(&clone); err != nil {
		s.logger.Errorf("Failed to clone transaction: %v", err)
		s.sendJSONError(w, "Failed to clone transaction", http.StatusInternalServerError)
		return
	}

	s.sendJSONSuccess(w, toTransactionDTO(clone))
}

// buildSearchFilter validates inputs and assembles a repository.TransactionFilter.
// Returns (filter, httpStatusOnError, errMessageOnError).
func buildSearchFilter(query, category, txType, dateFrom, dateTo string, amountMin, amountMax *float64) (repository.TransactionFilter, int, string) {
	var f repository.TransactionFilter
	f.Query = strings.TrimSpace(query)

	if category != "" && category != "all" {
		if !model.IsValidTransactionCategory(category) {
			return f, http.StatusBadRequest, "Invalid category"
		}
		f.Category = category
	}

	if txType != "" {
		if txType != string(model.TypeIncome) && txType != string(model.TypeExpense) {
			return f, http.StatusBadRequest, "Invalid transaction type"
		}
		f.Type = model.TransactionType(txType)
	}

	if dateFrom != "" {
		d, err := time.Parse(dateLayout, dateFrom)
		if err != nil {
			return f, http.StatusBadRequest, "Invalid dateFrom (expected YYYY-MM-DD)"
		}
		f.DateFrom = &d
	}
	if dateTo != "" {
		d, err := time.Parse(dateLayout, dateTo)
		if err != nil {
			return f, http.StatusBadRequest, "Invalid dateTo (expected YYYY-MM-DD)"
		}
		// Make dateTo inclusive to the end of the day.
		end := d.Add(24*time.Hour - time.Nanosecond)
		f.DateTo = &end
	}
	if f.DateFrom != nil && f.DateTo != nil && f.DateFrom.After(*f.DateTo) {
		return f, http.StatusBadRequest, "dateFrom must be on or before dateTo"
	}

	if amountMin != nil {
		if *amountMin < 0 {
			return f, http.StatusBadRequest, "amountMin must be >= 0"
		}
		f.AmountMin = amountMin
	}
	if amountMax != nil {
		if *amountMax < 0 {
			return f, http.StatusBadRequest, "amountMax must be >= 0"
		}
		f.AmountMax = amountMax
	}
	if f.AmountMin != nil && f.AmountMax != nil && *f.AmountMin > *f.AmountMax {
		return f, http.StatusBadRequest, "amountMin must be <= amountMax"
	}

	return f, 0, ""
}

// handleAPISearchTransactions returns transactions matching a filter set.
//
//	@Summary		Search transactions
//	@Description	Search a user's transactions by any combination of text, category, type, date range, amount range. Returns paginated results with a total count.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SearchTransactionsRequest	true	"Filter set"
//	@Success		200		{object}	SearchTransactionsResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/transactions/search [post]
func (s *Server) handleAPISearchTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SearchTransactionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	filter, code, msg := buildSearchFilter(req.Query, req.Category, req.Type, req.DateFrom, req.DateTo, req.AmountMin, req.AmountMax)
	if code != 0 {
		s.sendJSONError(w, msg, code)
		return
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}
	limit := req.Limit
	switch {
	case limit <= 0:
		limit = 50
	case limit > 200:
		limit = 200
	}

	txs, total, err := s.repositories.Transactions.SearchUserTransactionsFiltered(user.TgID, filter, offset, limit)
	if err != nil {
		s.logger.Errorf("Failed to search transactions: %v", err)
		s.sendJSONError(w, "Failed to search transactions", http.StatusInternalServerError)
		return
	}

	dtos := make([]TransactionDTO, len(txs))
	for i, tx := range txs {
		dtos[i] = toTransactionDTO(tx)
	}

	s.sendJSONSuccess(w, SearchTransactionsResponse{
		Transactions: dtos,
		Total:        total,
		Offset:       offset,
		Limit:        limit,
	})
}

// parseFloatQuery returns a pointer to the parsed value, or nil if absent.
func parseFloatQuery(r *http.Request, key string) (*float64, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s", key)
	}
	return &v, nil
}

// handleAPIExportTransactions streams a CSV of the user's transactions.
//
//	@Summary		Export transactions as CSV
//	@Description	Stream a CSV containing all transactions matching the optional filter set. Columns: date,type,category,amount,currency,description,created_at,updated_at.
//	@Tags			transactions
//	@Produce		text/csv
//	@Param			query		query		string	false	"Substring match on description (case-insensitive)"
//	@Param			category	query		string	false	"Category filter (\"all\" or empty disables it)"
//	@Param			type		query		string	false	"Transaction type: Income or Expense"
//	@Param			dateFrom	query		string	false	"Inclusive lower bound (YYYY-MM-DD)"
//	@Param			dateTo		query		string	false	"Inclusive upper bound (YYYY-MM-DD)"
//	@Param			amountMin	query		number	false	"Inclusive lower bound on amount"
//	@Param			amountMax	query		number	false	"Inclusive upper bound on amount"
//	@Success		200			{file}		file
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/transactions/export [get]
func (s *Server) handleAPIExportTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := client.GetUserFromContext(r.Context())
	if user == nil {
		s.sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	amountMin, err := parseFloatQuery(r, "amountMin")
	if err != nil {
		s.sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	amountMax, err := parseFloatQuery(r, "amountMax")
	if err != nil {
		s.sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter, code, msg := buildSearchFilter(q.Get("query"), q.Get("category"), q.Get("type"), q.Get("dateFrom"), q.Get("dateTo"), amountMin, amountMax)
	if code != 0 {
		s.sendJSONError(w, msg, code)
		return
	}

	// limit <= 0 disables LIMIT.
	txs, _, err := s.repositories.Transactions.SearchUserTransactionsFiltered(user.TgID, filter, 0, 0)
	if err != nil {
		s.logger.Errorf("Failed to export transactions: %v", err)
		s.sendJSONError(w, "Failed to export transactions", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("cashout_export_%s.csv", time.Now().Format(dateLayout))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))

	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{"tg_id", "date", "type", "category", "amount", "currency", "description", "created_at", "updated_at"}
	if err := cw.Write(header); err != nil {
		s.logger.Errorf("Failed to write CSV header: %v", err)
		return
	}
	for _, tx := range txs {
		row := []string{
			strconv.FormatInt(tx.TgID, 10),
			tx.Date.Format(dateLayout),
			string(tx.Type),
			string(tx.Category),
			strconv.FormatFloat(tx.Amount, 'f', 2, 64),
			string(tx.Currency),
			tx.Description,
			tx.CreatedAt.Format(time.RFC3339),
			tx.UpdatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(row); err != nil {
			s.logger.Errorf("Failed to write CSV row: %v", err)
			return
		}
	}
}
