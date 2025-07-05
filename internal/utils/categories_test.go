package utils

import (
	"cashout/internal/model"
	"testing"
)

func TestGetCategoryEmoji(t *testing.T) {
	tests := []struct {
		name     string
		category model.TransactionCategory
		want     string
	}{
		{
			name:     "Salary category",
			category: model.CategorySalary,
			want:     "💵",
		},
		{
			name:     "Car category",
			category: model.CategoryCar,
			want:     "🚗",
		},
		{
			name:     "Grocery category",
			category: model.CategoryGrocery,
			want:     "🛒",
		},
		{
			name:     "Health category",
			category: model.CategoryHealth,
			want:     "🏥",
		},
		{
			name:     "Pets category",
			category: model.CategoryPets,
			want:     "🐈",
		},
		{
			name:     "Invalid category returns default",
			category: model.TransactionCategory("InvalidCategory"),
			want:     "📌",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCategoryEmoji(tt.category); got != tt.want {
				t.Errorf("GetCategoryEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAnIncomeTransactionPrompt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "contains salary",
			input: "salary 3000 euro",
			want:  true,
		},
		{
			name:  "contains income",
			input: "monthly income received",
			want:  true,
		},
		{
			name:  "contains paycheck",
			input: "got my paycheck today",
			want:  true,
		},
		{
			name:  "contains refund",
			input: "amazon refund 50",
			want:  true,
		},
		{
			name:  "case insensitive",
			input: "SALARY payment",
			want:  true,
		},
		{
			name:  "word boundary check",
			input: "incomes for the month",
			want:  false, // "incomes" not "income"
		},
		{
			name:  "no income words",
			input: "coffee 3.50",
			want:  false,
		},
		{
			name:  "expense description",
			input: "groceries at supermarket",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAnIncomeTransactionPrompt(tt.input); got != tt.want {
				t.Errorf("IsAnIncomeTransactionPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
