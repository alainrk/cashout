package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"cashout/internal/model"
	"cashout/internal/utils"

	"github.com/sirupsen/logrus"
)

type LLM struct {
	APIKey   string
	Endpoint string
	Model    string
	Logger   *logrus.Logger
}

type ExtractedTransaction struct {
	Type        model.TransactionType
	Description string
	Amount      float64
	Category    string
	Date        time.Time
}

func (llm *LLM) ExtractTransaction(userText string, transactionType model.TransactionType) (ExtractedTransaction, error) {
	transaction := ExtractedTransaction{
		Type: transactionType,
	}

	tmpl := LLMExpensePromptTemplate
	if transactionType == model.TypeIncome {
		tmpl = LLMIncomePromptTemplate
	}

	// Generate prompt using the template
	prompt, err := GeneratePrompt(userText, tmpl)
	if err != nil {
		llm.Logger.Errorf("Error generating prompt: %v\n", err)
		return transaction, err
	}

	// Request payload
	requestBody, err := json.Marshal(map[string]any{
		"model": llm.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 250,
	})
	if err != nil {
		llm.Logger.Errorf("Error creating request: %v\n", err)
		return transaction, err
	}

	// Create request
	req, err := http.NewRequest("POST", llm.Endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		llm.Logger.Errorf("Error creating request: %v\n", err)
		return transaction, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.APIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		llm.Logger.Errorf("Error sending request: %v\n", err)
		return transaction, err
	}
	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		llm.Logger.Errorf("Error reading response: %v\n", err)
		return transaction, err
	}

	// Parse response
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		llm.Logger.Errorf("Error parsing response: %v\n", err)
		llm.Logger.Errorln("Raw response", body)
		return transaction, err
	}

	// Extract and print the message content
	var content string
	if choices, ok := result["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if message, ok := choice["message"].(map[string]any); ok {
				llm.Logger.Debugln("LLM Message", message)
				content = fmt.Sprintf("%v", message["content"])
			}
		}
	} else {
		llm.Logger.Errorln("Raw response", body)
		return transaction, fmt.Errorf("invalid response format")
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
	var transactionData map[string]any
	if err := json.Unmarshal([]byte(content), &transactionData); err != nil {
		llm.Logger.Errorln("Error parsing LLM response as JSON", err)
		return transaction, err
	}

	// Extract fields
	if description, ok := transactionData["description"].(string); ok {
		transaction.Description = description
	}

	if amount, ok := transactionData["amount"].(float64); ok {
		transaction.Amount = amount
	}

	if category, ok := transactionData["category"].(string); ok {
		transaction.Category = category
	}

	transaction.Date = time.Now()
	if date, ok := transactionData["date"].(string); ok {
		transaction.Date, err = utils.ParseDate(date)
		if err != nil {
			transaction.Date = time.Now()
		}
	}

	return transaction, nil
}
