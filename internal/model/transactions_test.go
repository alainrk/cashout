package model

import "testing"

func TestIsValidTransactionCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{
			name:     "valid income category",
			category: "Salary",
			want:     true,
		},
		{
			name:     "valid expense category",
			category: "Grocery",
			want:     true,
		},
		{
			name:     "valid other income",
			category: "OtherIncomes",
			want:     true,
		},
		{
			name:     "valid other expense",
			category: "OtherExpenses",
			want:     true,
		},
		{
			name:     "invalid category",
			category: "InvalidCategory",
			want:     false,
		},
		{
			name:     "empty category",
			category: "",
			want:     false,
		},
		{
			name:     "case sensitive check",
			category: "salary", // lowercase
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTransactionCategory(tt.category); got != tt.want {
				t.Errorf("IsValidTransactionCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionTypeValue(t *testing.T) {
	tests := []struct {
		name  string
		tType TransactionType
		want  string
	}{
		{
			name:  "income type",
			tType: TypeIncome,
			want:  "Income",
		},
		{
			name:  "expense type",
			tType: TypeExpense,
			want:  "Expense",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tType.Value()
			if err != nil {
				t.Errorf("TransactionType.Value() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("TransactionType.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
