package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) AddTransactionIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery
	msg := query.Message

	action := strings.Split(query.Data, ".")[2]

	switch action {
	case "income":
		user.Session.State = model.StateInsertingIncome
	case "expense":
		user.Session.State = model.StateInsertingExpense
	default:
		c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
		return nil
	}

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg.EditText(b, fmt.Sprintf("Sure! To add a new <b>%s</b>:\nTell me category, amount and description. You can also specify a date and change it later, today is default.\n\n<i>Examples:</i>\n<code>Irish Pub 3.4</code>\n<code>January salary 3k 10/01</code>", action), &gotgbot.EditMessageTextOpts{
		ParseMode: "HTML",
	})
	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Cancel", CallbackData: "transactions.cancel"}},
			},
		},
	})

	return ext.ContinueGroups
}

func (c *Client) addTransaction(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transactionType model.TransactionType

	switch user.Session.State {
	case model.StateInsertingIncome:
		transactionType = model.TypeIncome
	case model.StateInsertingExpense:
		transactionType = model.TypeExpense
	default:
		c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
		return nil
	}

	transaction, err := c.LLM.ExtractTransaction(ctx.Message.Text, transactionType)
	if err != nil {
		msg := "I'm sorry, I couldn't understand your transaction!"
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return err
	}

	if transaction.Amount == 0 {
		msg := "I'm sorry, I couldn't understand your transaction!"
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Store the transaction in the session
	user.Session.State = model.StateWaitingConfirm
	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, msg, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
				{
					{
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})
	if err != nil {
		c.Logger.Errorln("failed to send confirm message", err)
		return err
	}

	return nil
}

func (c *Client) EditTransactionIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	var transaction model.Transaction
	err = json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	query := ctx.CallbackQuery

	field := strings.Split(query.Data, ".")[2]

	var text string
	var opts *gotgbot.SendMessageOpts

	switch field {
	case "description":
		user.Session.State = model.StateEditingTransactionDescription

		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: "transactions.cancel",
				},
			},
		}
		opts = &gotgbot.SendMessageOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}}
		text = fmt.Sprintf("Enter a new description for the transaction:\n\nCurrent: %s ", transaction.Description)
	case "amount":
		user.Session.State = model.StateEditingTransactionAmount

		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: "transactions.cancel",
				},
			},
		}
		opts = &gotgbot.SendMessageOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}}
		text = fmt.Sprintf("Enter a new amount for the transaction:\n\nCurrent: %.2f€ ", transaction.Amount)
	case "date":
		user.Session.State = model.StateEditingTransactionDate
		text = "Add your date (e.g. dd mm, dd-mm, dd-mm-yyyy)."
		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: "transactions.cancel",
				},
			},
		}
		opts = &gotgbot.SendMessageOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}}
	case "category":
		user.Session.State = model.StateEditingTransactionCategory
		text = "Choose your category among the following ones."

		keyboard := [][]gotgbot.KeyboardButton{
			{{Text: "Cancel"}},
			{{Text: "Salary"}},
			{{Text: "OtherIncomes"}},
		}

		if transaction.Type == model.TypeExpense {
			keyboard = [][]gotgbot.KeyboardButton{
				{{Text: "Cancel"}},
				{{Text: "Car"}},
				{{Text: "Clothes"}},
				{{Text: "Grocery"}},
				{{Text: "House"}},
				{{Text: "Bills"}},
				{{Text: "Entertainment"}},
				{{Text: "Sport"}},
				{{Text: "EatingOut"}},
				{{Text: "Transport"}},
				{{Text: "Learning"}},
				{{Text: "Toiletry"}},
				{{Text: "Health"}},
				{{Text: "Tech"}},
				{{Text: "Gifts"}},
				{{Text: "Travel"}},
				{{Text: "OtherExpenses"}},
			}
		}

		opts = &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
				Keyboard:        keyboard,
				OneTimeKeyboard: true,
				IsPersistent:    false,
				ResizeKeyboard:  true,
			},
		}
	default:
		return fmt.Errorf("unknown field: %s", field)
	}

	c.CleanupKeyboard(b, ctx)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	b.SendMessage(ctx.EffectiveSender.ChatId, text, opts)

	return nil
}

func (c *Client) editTransactionDate(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	// Get date from DD-MM-YYYY to date
	date, err := utils.ParseDate(ctx.Message.Text)
	if err != nil {
		fmt.Printf("failed to parse date: %v\n", err)
		b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid date, please try again.", nil)
		return err
	}

	if date.After(time.Now()) {
		b.SendMessage(ctx.EffectiveSender.ChatId, "I don't support future dates, please try again.", nil)
		return fmt.Errorf("invalid date: %s", ctx.Message.Text)
	}

	transaction.Date = date

	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	m := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
				{
					{
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})

	return nil
}

func (c *Client) editTransactionAmount(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	// Parse new amount from message
	amountStr := strings.TrimSpace(ctx.Message.Text)
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	newAmount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Invalid amount. Please enter a valid number.",
			nil,
		)
		return err
	}

	if newAmount <= 0 {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Amount must be greater than zero.",
			nil,
		)
		return err
	}

	// Update the transaction
	transaction.Amount = newAmount

	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	m := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
				{
					{
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})

	return nil
}

func (c *Client) editTransactionDescription(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	transaction.Description = strings.TrimSpace(ctx.Message.Text)
	if transaction.Description == "" {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Description cannot be empty.",
			nil,
		)
		return err
	}

	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	m := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
				{
					{
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})

	return nil
}

func (c *Client) editTransactionCategory(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	if !model.IsValidTransactionCategory(ctx.Message.Text) {
		b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid category, please try again.", nil)
		return fmt.Errorf("invalid category: %s", ctx.Message.Text)
	}
	transaction.Category = model.TransactionCategory(ctx.Message.Text)

	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	m := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
				{
					{
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})

	return nil
}

// Confirm confirms the previous action after the user been prompted.
// This function handles both new transaction additions and edits confirmation.
func (c *Client) Confirm(b *gotgbot.Bot, ctx *ext.Context) error {
	requestCtx := ctx.Request.Context()
	if requestCtx == nil {
		requestCtx = context.Background()
	}
	startTime := time.Now()
	var operationErr error
	operationType := "unknown_confirm" // Default, will be refined

	// Determine operation type based on user session state
	// This is a simplification; a more robust state machine might be needed.
	// Assuming StateWaitingConfirm is for new transactions,
	// and other states like StateEditingTransaction* lead here for edits.
	// For now, we'll assume this Confirm is primarily for 'add' or a generic 'save'.
	// Specific 'edit' metrics will be in their respective functions.

	_, u := c.getUserFromContext(ctx)
	user, operationErr := c.authAndGetUser(u)
	if operationErr != nil {
		// Metric recording for auth failure or user retrieval failure before main logic
		defer func() {
			duration := time.Since(startTime).Seconds()
			status := "failure"
			// If user object is available, try to infer operation type from session state for more accurate metric
			if user.Session.State == model.StateWaitingConfirm { // Typically for new transactions
				operationType = "add"
			} else if strings.HasPrefix(string(user.Session.State), "StateEditingTransaction") { // States like StateEditingTransactionAmount
				operationType = "edit_confirm" // General confirm after an edit
			}

			c.TransactionOperationsCounter.Add(requestCtx, 1,
				metric.WithAttributes(
					attribute.String("operation.type", operationType),
					attribute.String("status", status),
					attribute.String("error.type", "authentication_user_retrieval"),
				),
			)
			c.TransactionOperationDuration.Record(requestCtx, duration,
				metric.WithAttributes(
					attribute.String("operation.type", operationType),
					attribute.String("status", status),
					attribute.String("error.type", "authentication_user_retrieval"),
				),
			)
		}()
		return operationErr
	}

	// Determine operation type more definitively now that we have user state
	if user.Session.State == model.StateWaitingConfirm { // Typically for new transactions
		operationType = "add"
	} else if strings.HasPrefix(string(user.Session.State), "StateEditingTransaction") || user.Session.State == model.StateTransactionEditingConfirm { // States like StateEditingTransactionAmount or confirming an edit
		operationType = "edit_confirm" // General confirm after an edit
	}


	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if operationErr != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("operation.type", operationType),
			attribute.String("status", status),
		}
		// Potentially add transaction.type (income/expense) if available and relevant
		// if transaction != nil && transaction.Type != "" {
		// 	attrs = append(attrs, attribute.String("transaction.type", string(transaction.Type)))
		// }

		c.TransactionOperationsCounter.Add(requestCtx, 1, metric.WithAttributes(attrs...))
		c.TransactionOperationDuration.Record(requestCtx, duration, metric.WithAttributes(attrs...))
	}()

	var transaction model.Transaction
	operationErr = json.Unmarshal([]byte(user.Session.Body), &transaction)
	if operationErr != nil {
		operationErr = fmt.Errorf("failed to extract transaction from the session: %w", operationErr)
		return operationErr
	}

	transaction.TgID = user.TgID
	transaction.Currency = model.CurrencyEUR

	// The actual database operation
	// If operationType is 'edit_confirm', this might be an update rather than add.
	// Assuming Repositories.Transactions.Add handles upsert or this path is only for new transactions.
	// For now, we assume this is the 'add' path as per original logic.
	// If this 'Confirm' is also for edits, the underlying DB call might change.
	// Let's stick to the current logic: Repositories.Transactions.Add
	// If it's an edit, the previous edit functions should have updated session.Body with the final transaction state.
	// And then this Confirm just saves it.
	// This means the operation is more like a "save_transaction"
	if operationType == "edit_confirm" {
		// If it's an edit confirmation, we should call an update method
		// For now, let's assume Repositories.Transactions.Add is smart enough or this is fine.
		// Ideally, there would be a c.Repositories.Transactions.Update() method.
		// This example will proceed with Add, but in a real scenario, this might need adjustment.
		// For the purpose of metrics, "edit_confirm" operation type is set.
		// Let's assume an Update method exists for edited transactions for more accurate DB representation.
		// operationErr = c.Repositories.Transactions.Update(transaction)
		// For now, sticking to Add as per existing code structure:
		operationErr = c.Repositories.Transactions.Add(transaction)

	} else { // "add"
		operationErr = c.Repositories.Transactions.Add(transaction)
	}


	if operationErr != nil {
		SendMessage(ctx, b, "There has been an error saving your transaction, please retry", nil)
		c.Logger.Errorln("failed to save transaction ("+operationType+")", operationErr)

		// Reset the state carefully
		// If it was an add, reset to inserting. If edit, maybe back to normal or specific edit state.
		// For now, mirroring existing reset logic.
		if transaction.Type == model.TypeExpense {
			user.Session.State = model.StateInsertingExpense
		} else {
			user.Session.State = model.StateInsertingIncome
		}
		// Not returning here, let defer handle metrics
		updateErr := c.Repositories.Users.Update(&user)
		if updateErr != nil {
			// Log this error but prioritize the original operationErr for return
			c.Logger.Errorf("failed to set user data to reset the state after save failure: %v", updateErr)
		}
		operationErr = fmt.Errorf("failed to save transaction (%s): %w", operationType, operationErr)
		return operationErr
	}

	user.Session.State = model.StateNormal
	user.Session.Body = ""

	updateErr := c.Repositories.Users.Update(&user)
	if updateErr != nil {
		// Log this error, but the main operation was successful.
		// Consider if this should mark the operation as partial failure for metrics.
		c.Logger.Errorf("failed to set user data after successful save: %v", updateErr)
		// For now, not changing operationErr, as main transaction succeeded.
		// operationErr = fmt.Errorf("failed to set user data post-save: %w", updateErr)
	}

	// Remove the keyboard from the previous message
	query := ctx.CallbackQuery
	msg := query.Message
	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{},
		},
	})

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}
	// Success message part of SendHomeKeyboard
	errSend := c.SendHomeKeyboard(b, ctx, fmt.Sprintf("%s Your transaction has been saved!", emoji))
	if errSend != nil {
		c.Logger.Errorf("Error sending home keyboard after confirm: %v", errSend)
		// The main operation was successful, but sending confirmation failed.
		// This could be a separate metric or log. For now, just log.
	}
	return operationErr // This will be nil if everything up to user.Session.Body clear was successful
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	c.CleanupKeyboard(b, ctx)
	c.SendHomeKeyboard(b, ctx, "Your operation has been canceled!\nWhat else can I do for you?\n\n/edit - Edit a transaction\n/delete - Delete a transaction\n/list - List your transactions\n/month Month Recap\n/year Year Recap")

	return ext.EndGroups
}
