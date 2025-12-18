// Package ai contains the AI prompts for the web dashboard
package ai

import (
	"bytes"
	"text/template"
)

// LLMExpensePromptTemplate is the LLM prompt template for expense transactions
const LLMExpensePromptTemplate = `You are a financial transaction parser. Your task is to analyze the input text and extract the following information:
- The category of the transaction
- The amount spent or received
- A brief description of the transaction

Format the result as a JSON object with the following structure:
{ "category": "Category", "amount": 12.34, "description": "Description" }

Available categories (use ONLY these):
"Car", "Clothes", "Grocery", "House", "Bills", "Entertainment", "Sport", "EatingOut", "Transport", "Learning", "Toiletry", "Health", "Tech", "Gifts", "Travel", "Pets", "OtherExpenses"

Follow these rules:
1. For category selection:
   - First try to find the category directly mentioned in the text (accounting for typos/synonyms)
   - If no category is directly mentioned, infer it from the description
   - If category cannot be determined, use "OtherExpenses"
2. For description:
   - Use the main item mentioned in the text
   - Capitalize the first letter of the description
   - If no item is mentioned, use text of the category
3. For amount:
   - Convert any amount to standard decimal notation with a period (not comma) as decimal separator
   - Return as a number (not a string) with at most 2 decimal places
   - If no amount is mentioned, use 0

Examples:
- "bread 5 euro an 20, grocery" → { "category": "Grocery", "amount": 5.2, "description": "Bread" }
- "pam 4.31 grocertw" → { "category": "Grocery", "amount": 4.31, "description": "Pam" }
- "car 25,30" → { "category": "Car", "amount": 25.3, "description": "Car" }
- "34 usd 23-04" → { "category": "OtherExpenses", "amount": 34, "description": "OtherExpenses" }
- "Great sea food 12 euro e 25" → { "category": "EatingOut", "amount": 12.25, "description": "Great see food" }

IMPORTANT: Respond with ONLY the JSON object but without markdown syntax. Your answer is plaintext being JSON to be parsed as it is, don't include the triple backticks syntax or anything similar.

User input:
{{.UserText}}
`

// LLMIncomePromptTemplate is the LLM prompt template for income transactions
const LLMIncomePromptTemplate = `You are a financial transaction parser. Your task is to analyze the input text and extract the following information:
- The category of the transaction
- The amount spent or received
- A brief description of the transaction

Format the result as a JSON object with the following structure:
{ "category": "Category", "amount": 12.34, "description": "Description" }

Available categories (use ONLY these):
"Salary", "OtherIncomes"

Follow these rules:
1. For category selection:
   - First try to find the category directly mentioned in the text (accounting for typos/synonyms)
   - If no category is directly mentioned, infer it from the description and prefer "Salary" only when the user use it or with a synonym in any language
   - If category cannot be determined, use "OtherIncomes", for example for "ticket restaurants", "refund amazon", etc.
2. For description:
   - Use the main item mentioned in the text
   - Capitalize the first letter of the description
   - If no item is mentioned, use text of the category
3. For amount:
   - Convert any amount to standard decimal notation with a period (not comma) as decimal separator
   - Return as a number (not a string) with at most 2 decimal places
   - If no amount is mentioned, use 0

Examples:
- "250k earned from job" → { "category": "Salary", "amount": 250000, "description": "From job" }
- "salayr 340 and 34 august" → { "category": "Salary", "amount": 340.34, "description": "August" }
- "ticket reastants 245 dollars" → { "category": "OtherIncomes", "amount": 245, "description": "Ticket restaurants" }
- "gained income 231 and 32 euro 03-04" → { "category": "Salary", "amount": 231.32, "description": "Salary" }

IMPORTANT: Respond with ONLY the JSON object but without markdown syntax. Your answer is plaintext being JSON to be parsed as it is, don't include the triple backticks syntax or anything similar.

User input:
{{.UserText}}
`

// LLMIntentClassificationPromptTemplate is the LLM prompt template for classifying user intent
const LLMIntentClassificationPromptTemplate = `You are an intent classifier for a financial tracking bot. Your task is to analyze the user's message and determine what action they want to perform.

Available intents (use ONLY these exact strings):
- "add_expense": User wants to add/record an expense (spending money)
- "add_income": User wants to add/record income (receiving money)
- "edit": User wants to edit/modify/change an existing transaction
- "delete": User wants to delete/remove an existing transaction
- "search": User wants to search/find transactions
- "list": User wants to list/view/see all or recent transactions
- "week_recap": User wants to see a weekly summary/recap
- "month_recap": User wants to see a monthly summary/recap
- "year_recap": User wants to see a yearly summary/recap
- "export": User wants to export/download transactions (CSV, file)
- "unknown": Cannot determine the intent or it doesn't match any of the above

Classification rules:
1. If the message contains an amount (numbers with currency context), classify as "add_expense" unless income-related words are present
2. Income-related words: salary, wage, income, earned, received, got paid, paycheck, bonus, refund, reimbursement
3. Expense-related context: bought, spent, paid, cost, purchase
4. Edit-related words: edit, modify, change, update, fix, correct
5. Delete-related words: delete, remove, cancel, undo
6. Search-related words: search, find, look for, where is, show me
7. List-related words: list, show all, view, display, transactions, history
8. Recap-related words: recap, summary, overview, total, how much
9. Export-related words: export, download, CSV, file, backup
10. If the message is a greeting, question about the bot, or unrelated to finance, use "unknown"

Format the result as a JSON object:
{ "intent": "intent_name", "confidence": 0.95 }

Where confidence is a value between 0 and 1 indicating how confident you are in the classification.

IMPORTANT: Respond with ONLY the JSON object without markdown syntax. Your answer is plaintext JSON to be parsed directly.

User input:
{{.UserText}}
`

// GeneratePrompt creates the complete prompt by filling in the template with user input
func GeneratePrompt(userText string, promptTemplate string) (string, error) {
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		UserText string
	}{
		UserText: userText,
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
