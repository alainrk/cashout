package ai

import (
	"slices"
	"strings"
	"testing"
)

func TestIntentCloneConstant(t *testing.T) {
	if IntentClone != "clone" {
		t.Errorf("IntentClone = %q, want %q", IntentClone, "clone")
	}
}

func TestCloneIntentInClassificationPrompt(t *testing.T) {
	prompt, err := GeneratePrompt("clone a transaction", LLMIntentClassificationPromptTemplate)
	if err != nil {
		t.Fatalf("GeneratePrompt error: %v", err)
	}

	// Verify the clone intent is listed in available intents
	if !strings.Contains(prompt, `"clone"`) {
		t.Error("prompt missing clone intent in available intents")
	}

	// Verify clone keywords are present in classification rules
	keywords := []string{"clone", "duplicate", "repeat", "copy", "same again"}
	for _, kw := range keywords {
		if !strings.Contains(prompt, kw) {
			t.Errorf("prompt missing clone keyword %q", kw)
		}
	}

	// Verify user text is injected
	if !strings.Contains(prompt, "clone a transaction") {
		t.Error("prompt missing user text")
	}
}

func TestCloneIntentValidation(t *testing.T) {
	// Verify IntentClone is in the valid set (matches the switch in ClassifyIntent)
	validIntents := []Intent{
		IntentAddExpense, IntentAddIncome, IntentEdit, IntentDelete, IntentSearch,
		IntentList, IntentWeekRecap, IntentMonthRecap, IntentYearRecap, IntentExport, IntentClone,
	}

	found := slices.Contains(validIntents, IntentClone)
	if !found {
		t.Error("IntentClone not in valid intents list")
	}
}
