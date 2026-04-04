package client

import (
	"cashout/internal/model"
	"strings"
	"testing"
	"time"
)

func TestFormatCloneRecentExpenses_Basic(t *testing.T) {
	txns := []model.Transaction{
		{
			ID:          1,
			Description: "Coffee Shop",
			Amount:      3.50,
			Category:    model.CategoryEatingOut,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Description: "Lidl",
			Amount:      45.00,
			Category:    model.CategoryGrocery,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatCloneRecentExpenses(txns, 0, 2)

	if !strings.Contains(result, "📋 <b>Clone Transaction</b>") {
		t.Error("missing header")
	}
	if !strings.Contains(result, "Recent expenses — 1–2 of 2") {
		t.Errorf("missing count, got:\n%s", result)
	}
	if !strings.Contains(result, "1. 🍽️ Coffee Shop · €3.50 · 01/04/2026") {
		t.Errorf("missing first transaction, got:\n%s", result)
	}
	if !strings.Contains(result, "2. 🛒 Lidl · €45.00 · 30/03/2026") {
		t.Errorf("missing second transaction, got:\n%s", result)
	}
	if !strings.Contains(result, "Tap a number to clone it with today's date.") {
		t.Error("missing tap instruction")
	}
}

func TestFormatCloneRecentExpenses_Offset(t *testing.T) {
	txns := []model.Transaction{
		{
			ID:          5,
			Description: "Bus Ticket",
			Amount:      2.00,
			Category:    model.CategoryTransport,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatCloneRecentExpenses(txns, 10, 15)
	if !strings.Contains(result, "Recent expenses — 11–11 of 15") {
		t.Errorf("offset counting wrong, got:\n%s", result)
	}
}

func TestCreateCloneRecentKeyboard_NumberedButtons(t *testing.T) {
	txns := []model.Transaction{
		{ID: 10}, {ID: 20}, {ID: 30},
	}

	kb := createCloneRecentKeyboard(txns, 0, 10, 3)

	// First row: 3 numbered buttons
	if len(kb[0]) != 3 {
		t.Fatalf("expected 3 numbered buttons, got %d", len(kb[0]))
	}
	if kb[0][0].Text != "1" || kb[0][0].CallbackData != "clone.select.10" {
		t.Errorf("button 1 wrong: %+v", kb[0][0])
	}
	if kb[0][1].Text != "2" || kb[0][1].CallbackData != "clone.select.20" {
		t.Errorf("button 2 wrong: %+v", kb[0][1])
	}
	if kb[0][2].Text != "3" || kb[0][2].CallbackData != "clone.select.30" {
		t.Errorf("button 3 wrong: %+v", kb[0][2])
	}
}

func TestCreateCloneRecentKeyboard_NumberedButtonsRowsOf5(t *testing.T) {
	txns := make([]model.Transaction, 7)
	for i := range txns {
		txns[i].ID = int64(i + 1)
	}

	kb := createCloneRecentKeyboard(txns, 0, 10, 7)

	// Row 0: 5 buttons, Row 1: 2 buttons, Row 2: SearchMore+Cancel
	if len(kb[0]) != 5 {
		t.Errorf("first row should have 5 buttons, got %d", len(kb[0]))
	}
	if len(kb[1]) != 2 {
		t.Errorf("second row should have 2 buttons, got %d", len(kb[1]))
	}
}

func TestCreateCloneRecentKeyboard_SearchMoreAndCancel(t *testing.T) {
	txns := []model.Transaction{{ID: 1}}
	kb := createCloneRecentKeyboard(txns, 0, 10, 1)

	// Last row should be Search More + Cancel
	lastRow := kb[len(kb)-1]
	if len(lastRow) != 2 {
		t.Fatalf("last row should have 2 buttons, got %d", len(lastRow))
	}
	if lastRow[0].CallbackData != "clone.searchmore" {
		t.Errorf("expected clone.searchmore, got %s", lastRow[0].CallbackData)
	}
	if lastRow[1].CallbackData != "transactions.cancel" {
		t.Errorf("expected transactions.cancel, got %s", lastRow[1].CallbackData)
	}
}

func TestCreateCloneRecentKeyboard_SinglePage_NoNavigation(t *testing.T) {
	txns := []model.Transaction{{ID: 1}, {ID: 2}}
	kb := createCloneRecentKeyboard(txns, 0, 10, 2)

	// Should be: numbered row + search/cancel row (no navigation)
	if len(kb) != 2 {
		t.Errorf("single page should have 2 rows (numbers + actions), got %d", len(kb))
	}
}

func TestCreateCloneRecentKeyboard_MultiPage_HasNavigation(t *testing.T) {
	txns := make([]model.Transaction, 10)
	for i := range txns {
		txns[i].ID = int64(i + 1)
	}

	// First page of 25 total
	kb := createCloneRecentKeyboard(txns, 0, 10, 25)

	// Should have: 2 number rows (5+5) + nav row + search/cancel row
	if len(kb) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(kb))
	}

	// Nav row (third row): page indicator + no Next (first page)
	navRow := kb[2]
	foundPrev := false
	foundIndicator := false
	for _, btn := range navRow {
		if btn.Text == "⬅️ Previous" {
			foundPrev = true
		}
		if btn.Text == "1/3" {
			foundIndicator = true
		}
	}
	if !foundPrev {
		t.Error("first page should have Previous button (shows older)")
	}
	if !foundIndicator {
		t.Error("missing page indicator 1/3")
	}
}

func TestCreateCloneRecentKeyboard_MiddlePage_HasBothNav(t *testing.T) {
	txns := make([]model.Transaction, 10)
	for i := range txns {
		txns[i].ID = int64(i + 1)
	}

	// Middle page (offset=10, total=30)
	kb := createCloneRecentKeyboard(txns, 10, 10, 30)

	var navRow []string
	for _, row := range kb {
		for _, btn := range row {
			if btn.Text == "⬅️ Previous" || btn.Text == "Next ➡️" || strings.Contains(btn.Text, "/") {
				navRow = append(navRow, btn.Text)
			}
		}
	}

	if len(navRow) != 3 {
		t.Errorf("middle page should have Previous, indicator, Next, got %v", navRow)
	}
}

func TestFormatCloneSearchResults_WithQuery(t *testing.T) {
	txns := []model.Transaction{
		{
			ID:          1,
			Description: "Coffee at Starbucks",
			Amount:      4.50,
			Category:    model.CategoryEatingOut,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatCloneSearchResults(txns, "coffee", "all", 0, 1)

	if !strings.Contains(result, "📋 <b>Clone Transaction</b>") {
		t.Error("missing header")
	}
	if !strings.Contains(result, `Query: "coffee"`) {
		t.Errorf("missing query, got:\n%s", result)
	}
	if !strings.Contains(result, "<b>Coffee</b>") {
		t.Errorf("search term should be bolded, got:\n%s", result)
	}
	if !strings.Contains(result, "-€4.50") {
		t.Errorf("missing expense sign, got:\n%s", result)
	}
}

func TestFormatCloneSearchResults_WithCategory(t *testing.T) {
	txns := []model.Transaction{
		{
			ID:          1,
			Description: "Lidl",
			Amount:      30.00,
			Category:    model.CategoryGrocery,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatCloneSearchResults(txns, "%", "Grocery", 0, 1)

	if !strings.Contains(result, "🛒 Grocery") {
		t.Errorf("missing category filter, got:\n%s", result)
	}
	// With "%" query, should not show query text
	if strings.Contains(result, `Query:`) {
		t.Errorf("should not show query for wildcard, got:\n%s", result)
	}
}

func TestFormatCloneSearchResults_IncomeSign(t *testing.T) {
	txns := []model.Transaction{
		{
			ID:          1,
			Description: "March Salary",
			Amount:      3000.00,
			Category:    model.CategorySalary,
			Type:        model.TypeIncome,
			Date:        time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatCloneSearchResults(txns, "%", "all", 0, 1)

	if !strings.Contains(result, "+€3000.00") {
		t.Errorf("income should have + sign, got:\n%s", result)
	}
}

func TestCreateCloneSearchKeyboard_NumberedButtons(t *testing.T) {
	txns := []model.Transaction{{ID: 100}, {ID: 200}}

	kb := createCloneSearchKeyboard(txns, "all", "coffee", 0, 10, 2)

	if kb[0][0].CallbackData != "clone.search.select.100" {
		t.Errorf("wrong callback: %s", kb[0][0].CallbackData)
	}
	if kb[0][1].CallbackData != "clone.search.select.200" {
		t.Errorf("wrong callback: %s", kb[0][1].CallbackData)
	}
}

func TestCreateCloneSearchKeyboard_NewSearchAndHome(t *testing.T) {
	txns := []model.Transaction{{ID: 1}}

	kb := createCloneSearchKeyboard(txns, "all", "test", 0, 10, 1)

	lastRow := kb[len(kb)-1]
	if lastRow[0].CallbackData != "clone.search.new" {
		t.Errorf("expected clone.search.new, got %s", lastRow[0].CallbackData)
	}
	if lastRow[1].CallbackData != "clone.search.home" {
		t.Errorf("expected clone.search.home, got %s", lastRow[1].CallbackData)
	}
}

func TestCreateCloneSearchKeyboard_Pagination(t *testing.T) {
	txns := make([]model.Transaction, 10)
	for i := range txns {
		txns[i].ID = int64(i + 1)
	}

	// Middle page
	kb := createCloneSearchKeyboard(txns, "Grocery", "lidl", 10, 10, 30)

	// Find navigation row
	var navTexts []string
	for _, row := range kb {
		for _, btn := range row {
			if btn.Text == "⬅️ Previous" || btn.Text == "Next ➡️" || strings.Contains(btn.Text, "/") {
				navTexts = append(navTexts, btn.Text)
				// Check category and query are in callback data
				if btn.Text == "⬅️ Previous" {
					if !strings.Contains(btn.CallbackData, "Grocery") || !strings.Contains(btn.CallbackData, "lidl") {
						t.Errorf("prev button missing category/query: %s", btn.CallbackData)
					}
				}
			}
		}
	}

	if len(navTexts) != 3 {
		t.Errorf("expected 3 nav buttons, got %v", navTexts)
	}
}

func TestCreateCloneSearchKeyboard_SinglePage_NoNavigation(t *testing.T) {
	txns := []model.Transaction{{ID: 1}}

	kb := createCloneSearchKeyboard(txns, "all", "test", 0, 10, 1)

	// Should be: numbered row + actions row (no navigation)
	if len(kb) != 2 {
		t.Errorf("single page should have 2 rows, got %d", len(kb))
	}
}
