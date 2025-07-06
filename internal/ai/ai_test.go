package ai

import (
	"strings"
	"testing"
)

func TestGeneratePrompt(t *testing.T) {
	tests := []struct {
		name           string
		userText       string
		promptTemplate string
		wantContains   []string
		wantErr        bool
	}{
		{
			name:           "expense prompt generation",
			userText:       "coffee 3.50",
			promptTemplate: LLMExpensePromptTemplate,
			wantContains: []string{
				"coffee 3.50",
				"financial transaction parser",
				"OtherExpenses",
			},
			wantErr: false,
		},
		{
			name:           "income prompt generation",
			userText:       "salary 3000",
			promptTemplate: LLMIncomePromptTemplate,
			wantContains: []string{
				"salary 3000",
				"Salary",
				"OtherIncomes",
			},
			wantErr: false,
		},
		{
			name:           "empty user text",
			userText:       "",
			promptTemplate: LLMExpensePromptTemplate,
			wantContains:   []string{"User input:"},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePrompt(tt.userText, tt.promptTemplate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("GeneratePrompt() result doesn't contain %q", want)
				}
			}
		})
	}
}
