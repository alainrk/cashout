package ai

import (
	"cashout/internal/model"
	"encoding/json"
	"testing"
	"time"
)

func TestExtractedTransactionJSON(t *testing.T) {
	// Test that ExtractedTransaction can be properly marshaled/unmarshaled
	transaction := ExtractedTransaction{
		Type:        model.TypeExpense,
		Description: "Coffee at Starbucks",
		Amount:      3.50,
		Category:    "EatingOut",
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	// Marshal to JSON
	data, err := json.Marshal(transaction)
	if err != nil {
		t.Fatalf("Failed to marshal ExtractedTransaction: %v", err)
	}

	// Unmarshal back
	var unmarshaled ExtractedTransaction
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ExtractedTransaction: %v", err)
	}

	// Compare
	if unmarshaled.Type != transaction.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaled.Type, transaction.Type)
	}
	if unmarshaled.Description != transaction.Description {
		t.Errorf("Description mismatch: got %v, want %v", unmarshaled.Description, transaction.Description)
	}
	if unmarshaled.Amount != transaction.Amount {
		t.Errorf("Amount mismatch: got %v, want %v", unmarshaled.Amount, transaction.Amount)
	}
	if unmarshaled.Category != transaction.Category {
		t.Errorf("Category mismatch: got %v, want %v", unmarshaled.Category, transaction.Category)
	}
}
