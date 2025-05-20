package ai

import (
	"bytes"
	"cashout/internal/model"
	"cashout/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log" // Using standard log for init error, as LLM logger might not be ready
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	aiAPICallsCounter metric.Int64Counter
	aiAPICallDuration metric.Float64Histogram
)

func init() {
	meter := otel.Meter("cashout/ai")
	var err error
	aiAPICallsCounter, err = meter.Int64Counter(
		"ai.api.calls.total",
		metric.WithDescription("Counts the number of AI API calls."),
	)
	if err != nil {
		log.Printf("Error initializing aiAPICallsCounter: %v\n", err)
		// Depending on policy, could panic. For now, we log and the counter will be nil.
	}

	aiAPICallDuration, err = meter.Float64Histogram(
		"ai.api.call.duration.seconds",
		metric.WithDescription("Measures the duration of AI API calls in seconds."),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Printf("Error initializing aiAPICallDuration: %v\n", err)
	}
}

type LLM struct {
	APIKey   string
	Endpoint string
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
	// Assuming context.Background() for now as the function signature does not include context.
	// Ideally, context should be passed down from the caller.
	requestCtx := context.Background()
	startTime := time.Now()
	var operationErr error // Use a specific variable for the operation's error.

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if operationErr != nil {
			status = "failure"
		}

		if aiAPICallsCounter != nil {
			aiAPICallsCounter.Add(requestCtx, 1, metric.WithAttributes(attribute.String("status", status)))
		} else {
			llm.Logger.Warn("aiAPICallsCounter is not initialized. Skipping metric.")
		}

		if aiAPICallDuration != nil {
			aiAPICallDuration.Record(requestCtx, duration, metric.WithAttributes(attribute.String("status", status)))
		} else {
			llm.Logger.Warn("aiAPICallDuration is not initialized. Skipping metric.")
		}
	}()

	transaction := ExtractedTransaction{
		Type: transactionType,
	}

	tmpl := LLMExpensePromptTemplate
	if transactionType == model.TypeIncome {
		tmpl = LLMIncomePromptTemplate
	}

	// Generate prompt using the template
	prompt, genPromptErr := GeneratePrompt(userText, tmpl)
	if genPromptErr != nil {
		llm.Logger.Errorf("Error generating prompt: %v\n", genPromptErr)
		operationErr = genPromptErr
		return transaction, operationErr
	}

	// Request payload
	requestPayload := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 250,
	}
	requestBody, marshalErr := json.Marshal(requestPayload)
	if marshalErr != nil {
		llm.Logger.Errorf("Error creating request JSON: %v\n", marshalErr)
		operationErr = marshalErr
		return transaction, operationErr
	}

	// Create request
	req, newReqErr := http.NewRequestWithContext(requestCtx, "POST", llm.Endpoint, bytes.NewBuffer(requestBody))
	if newReqErr != nil {
		llm.Logger.Errorf("Error creating HTTP request: %v\n", newReqErr)
		operationErr = newReqErr
		return transaction, operationErr
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.APIKey)

	// Send request
	client := &http.Client{}
	resp, doReqErr := client.Do(req)
	if doReqErr != nil {
		llm.Logger.Errorf("Error sending HTTP request: %v\n", doReqErr)
		operationErr = doReqErr
		return transaction, operationErr
	}
	defer resp.Body.Close()

	// Read response
	body, readBodyErr := io.ReadAll(resp.Body)
	if readBodyErr != nil {
		llm.Logger.Errorf("Error reading response body: %v\n", readBodyErr)
		operationErr = readBodyErr
		return transaction, operationErr
	}

	if resp.StatusCode != http.StatusOK {
		llm.Logger.Errorf("AI API request failed with status %d: %s", resp.StatusCode, string(body))
		operationErr = fmt.Errorf("AI API request failed with status %d", resp.StatusCode)
		return transaction, operationErr
	}

	// Parse response
	var result map[string]interface{}
	if unmarshalRespErr := json.Unmarshal(body, &result); unmarshalRespErr != nil {
		llm.Logger.Errorf("Error parsing API response JSON: %v\n", unmarshalRespErr)
		llm.Logger.Errorln("Raw response body:", string(body))
		operationErr = unmarshalRespErr
		return transaction, operationErr
	}

	// Extract and print the message content
	var content string
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				llm.Logger.Debugln("LLM Message", message)
				if msgContent, ok := message["content"].(string); ok {
					content = msgContent
				} else {
					llm.Logger.Errorln("LLM message content is not a string:", message["content"])
					operationErr = fmt.Errorf("LLM message content is not a string")
					return transaction, operationErr
				}
			} else {
				llm.Logger.Errorln("Invalid 'message' format in LLM choice:", choice)
				operationErr = fmt.Errorf("invalid 'message' format in LLM choice")
				return transaction, operationErr
			}
		} else {
			llm.Logger.Errorln("Invalid 'choice' format in LLM response:", choices[0])
			operationErr = fmt.Errorf("invalid 'choice' format in LLM response")
			return transaction, operationErr
		}
	} else {
		llm.Logger.Errorln("No 'choices' in LLM response. Raw response body:", string(body))
		operationErr = fmt.Errorf("no 'choices' in LLM response")
		return transaction, operationErr
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
	var transactionData map[string]interface{}
	if parseContentErr := json.Unmarshal([]byte(content), &transactionData); parseContentErr != nil {
		llm.Logger.Errorln("Error parsing LLM content JSON:", parseContentErr, "Content:", content)
		operationErr = parseContentErr
		return transaction, operationErr
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

	var dateParseErr error
	transaction.Date = time.Now() // Default to now
	if dateStr, ok := transactionData["date"].(string); ok {
		parsedDate, errDate := utils.ParseDate(dateStr)
		if errDate != nil {
			llm.Logger.Warnf("Error parsing date from LLM '%s': %v. Defaulting to now.", dateStr, errDate)
			// Keep default date (time.Now()), do not set operationErr for this minor issue.
			dateParseErr = errDate // Store for potential logging but not critical for operation status
		} else {
			transaction.Date = parsedDate
		}
	}
	_ = dateParseErr // Avoid unused variable error if not logging it further here.

	return transaction, operationErr // operationErr will be nil if everything succeeded
}
