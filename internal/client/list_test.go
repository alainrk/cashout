package client

import (
	"cashout/internal/model"
	"strings"
	"testing"
	"time"
)

func TestFormatTransactions_Empty(t *testing.T) {
	result := formatTransactions(2026, 2, nil, 0, 0, "all")
	expected := "No transactions found for February 2026"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestFormatTransactions_EmptyWithCategory(t *testing.T) {
	result := formatTransactions(2026, 2, nil, 0, 0, "Grocery")
	if !strings.Contains(result, "No transactions found for February 2026") {
		t.Error("missing base message")
	}
	if !strings.Contains(result, "🛒 Grocery") {
		t.Error("missing category in empty message")
	}
}

func TestFormatTransactions_CompactSingleLine(t *testing.T) {
	txns := []model.Transaction{
		{
			Description: "Grocery Shopping",
			Amount:      45.00,
			Category:    model.CategoryGrocery,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Description: "January Salary",
			Amount:      3000.00,
			Category:    model.CategorySalary,
			Type:        model.TypeIncome,
			Date:        time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatTransactions(2026, 2, txns, 0, 2, "all")

	if !strings.Contains(result, "<b>February 2026</b>") {
		t.Error("missing month/year header")
	}
	if !strings.Contains(result, "Showing 1–2 of 2") {
		t.Error("missing showing count")
	}

	// Compact: emoji <b>desc</b> · sign€amount · DD/MM
	if !strings.Contains(result, "🛒 <b>Grocery Shopping</b> · -€45.00 · 10/02") {
		t.Errorf("missing compact grocery expense line, got:\n%s", result)
	}
	if !strings.Contains(result, "💵 <b>January Salary</b> · +€3000.00 · 01/02") {
		t.Errorf("missing compact salary income line, got:\n%s", result)
	}
}

func TestFormatTransactions_WithCategoryHeader(t *testing.T) {
	txns := []model.Transaction{
		{
			Description: "Lidl",
			Amount:      25.00,
			Category:    model.CategoryGrocery,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatTransactions(2026, 2, txns, 0, 1, "Grocery")
	if !strings.Contains(result, "🛒 Grocery") {
		t.Errorf("missing category in header, got:\n%s", result)
	}
}

func TestFormatTransactions_Offset(t *testing.T) {
	txns := []model.Transaction{
		{
			Description: "Test",
			Amount:      10.00,
			Category:    model.CategoryBills,
			Type:        model.TypeExpense,
			Date:        time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	result := formatTransactions(2026, 2, txns, 10, 15, "all")
	if !strings.Contains(result, "Showing 11–11 of 15") {
		t.Errorf("offset counting wrong, got:\n%s", result)
	}
}

func TestCreatePaginationKeyboard_PageIndicator(t *testing.T) {
	kb := createPaginationKeyboard(2026, 2, 10, 10, 30, "all")

	navRow := kb[0]
	if len(navRow) != 3 {
		t.Fatalf("expected 3 buttons in nav row, got %d", len(navRow))
	}
	if navRow[0].Text != "⬅️ Previous" {
		t.Errorf("expected Previous button, got %q", navRow[0].Text)
	}
	if navRow[1].Text != "2/3" {
		t.Errorf("expected page indicator '2/3', got %q", navRow[1].Text)
	}
	if navRow[2].Text != "Next ➡️" {
		t.Errorf("expected Next button, got %q", navRow[2].Text)
	}
}

func TestCreatePaginationKeyboard_CategoryInCallbackData(t *testing.T) {
	kb := createPaginationKeyboard(2026, 2, 10, 10, 30, "Grocery")

	navRow := kb[0]
	if !strings.Contains(navRow[0].CallbackData, "Grocery") {
		t.Errorf("prev button missing category, got %q", navRow[0].CallbackData)
	}
	if !strings.Contains(navRow[2].CallbackData, "Grocery") {
		t.Errorf("next button missing category, got %q", navRow[2].CallbackData)
	}

	backRow := kb[1]
	if !strings.Contains(backRow[0].CallbackData, "Grocery") {
		t.Errorf("back button missing category, got %q", backRow[0].CallbackData)
	}
}

func TestCreatePaginationKeyboard_FirstPage(t *testing.T) {
	kb := createPaginationKeyboard(2026, 2, 0, 10, 25, "all")

	navRow := kb[0]
	if len(navRow) != 2 {
		t.Fatalf("expected 2 buttons on first page, got %d", len(navRow))
	}
	if navRow[0].Text != "1/3" {
		t.Errorf("expected page indicator '1/3', got %q", navRow[0].Text)
	}
}

func TestCreatePaginationKeyboard_LastPage(t *testing.T) {
	kb := createPaginationKeyboard(2026, 2, 20, 10, 25, "all")

	navRow := kb[0]
	if len(navRow) != 2 {
		t.Fatalf("expected 2 buttons on last page, got %d", len(navRow))
	}
	if navRow[0].Text != "⬅️ Previous" {
		t.Errorf("expected Previous button, got %q", navRow[0].Text)
	}
	if navRow[1].Text != "3/3" {
		t.Errorf("expected page indicator '3/3', got %q", navRow[1].Text)
	}
}

func TestCreatePaginationKeyboard_SinglePage(t *testing.T) {
	kb := createPaginationKeyboard(2026, 2, 0, 10, 5, "all")

	navRow := kb[0]
	if len(navRow) != 1 {
		t.Fatalf("expected 1 button on single page, got %d", len(navRow))
	}
	if navRow[0].Text != "1/1" {
		t.Errorf("expected page indicator '1/1', got %q", navRow[0].Text)
	}

	if kb[1][0].Text != "🔙 Back to Months" {
		t.Errorf("expected back button, got %q", kb[1][0].Text)
	}
}

func TestCreatePaginationKeyboard_ZeroTotal(t *testing.T) {
	kb := createPaginationKeyboard(2026, 2, 0, 20, 0, "Pets")

	// Should only have the "Back to Months" row, no nav row
	if len(kb) != 1 {
		t.Fatalf("expected 1 row (back only), got %d", len(kb))
	}
	if kb[0][0].Text != "🔙 Back to Months" {
		t.Errorf("expected back button, got %q", kb[0][0].Text)
	}
}
