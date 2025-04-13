package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"happypoor/internal/model"
	"io"
	"net/http"
)

type LLM struct {
	APIKey   string
	Endpoint string
}

type ExtractedTransaction struct {
	Type        model.TransactionType
	Description string
	Amount      float64
	Category    string
}

func (llm *LLM) ExtractTransaction(userText string, transactionType model.TransactionType) (ExtractedTransaction, error) {
	expense := ExtractedTransaction{
		Type: transactionType,
	}

	tmpl := LLMExpensePromptTemplate
	if transactionType == model.TypeIncome {
		tmpl = LLMIncomePromptTemplate
	}

	// Generate prompt using the template
	prompt, err := GeneratePrompt(userText, tmpl)
	if err != nil {
		fmt.Printf("Error generating prompt: %v\n", err)
		return expense, err
	}

	// Request payload
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 1000,
	})
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return expense, err
	}

	// Create request
	req, err := http.NewRequest("POST", llm.Endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return expense, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.APIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return expense, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return expense, err
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		fmt.Println("Raw response:", string(body))
		return expense, err
	}

	// Extract and print the message content
	var content string
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				fmt.Printf("%+v	", message)
				content = fmt.Sprintf("%v", message["content"])
			}
		}
	} else {
		fmt.Println("Raw response:", string(body))
		return expense, fmt.Errorf("invalid response format")
	}

	// Sometimes the llm returns the ```json``` markdown format, despite being asked no to, so we need to clean it up
	jsonStart := 0
	jsonEnd := len(content)
	// Start parsing char by char until a "{" is found
	for i, char := range content {
		if char == '{' {
			jsonStart = i
			break
		}
	}
	// Starting from the end do the same until a "}" is found
	for i := len(content) - 1; i >= 0; i-- {
		if content[i] == '}' {
			jsonEnd = i + 1
			break
		}
	}
	// Remove the markdown
	content = content[jsonStart:jsonEnd]

	// ExtractExpense from the LLM Response text
	// Parse the LLM JSON response
	var expenseData map[string]interface{}
	if err := json.Unmarshal([]byte(content), &expenseData); err != nil {
		fmt.Printf("Error parsing LLM response as JSON: %v\n", err)
		return expense, err
	}

	// Extract fields
	if description, ok := expenseData["description"].(string); ok {
		expense.Description = description
	}

	if amount, ok := expenseData["amount"].(float64); ok {
		expense.Amount = amount
	}

	if category, ok := expenseData["category"].(string); ok {
		expense.Category = category
	}

	return expense, nil
}
