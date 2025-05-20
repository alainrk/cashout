package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// EditTransactions handles the /edit command
func (c *Client) EditTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user state to normal
	user.Session.State = model.StateNormal

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Show first page of editable transactions
	return c.showEditableTransactionPage(b, ctx, user, 0)
}

// EditTransactionPage handles pagination in the transaction editing interface
func (c *Client) EditTransactionPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: edit.page.OFFSET)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	offset, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}

	return c.showEditableTransactionPage(b, ctx, user, offset)
}

// EditTransactionSelect handles selection of a transaction to edit
func (c *Client) EditTransactionSelect(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: edit.select.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	// Get the transaction
	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Verify ownership
	if transaction.TgID != user.TgID {
		_, _, err = ctx.CallbackQuery.Message.EditText(
			b,
			"⚠️ This transaction doesn't belong to you.",
			&gotgbot.EditMessageTextOpts{},
		)
		return err
	}

	// Store transaction in session for later use
	user.Session.Body = fmt.Sprintf("%d", transactionID)
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Show edit options
	return c.showEditOptions(b, ctx, transaction)
}

// EditTransactionField handles editing a specific field of a transaction
func (c *Client) EditTransactionField(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Get transaction ID from session
	transactionID, err := strconv.ParseInt(user.Session.Body, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID in session: %v", err)
	}

	// Get the transaction
	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	query := ctx.CallbackQuery
	// Parse callback data (format: edit.field.FIELD_NAME)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	field := parts[2]

	switch field {
	case "description":
		return c.editTopLevelTransactionDescription(b, ctx, transaction)
	case "category":
		return c.editTopLevelTransactionCategory(b, ctx, transaction)
	case "amount":
		return c.editTopLevelTransactionAmount(b, ctx, transaction)
	case "date":
		return c.editTopLevelTransactionDate(b, ctx, transaction)
	default:
		return fmt.Errorf("invalid field: %s", field)
	}
}

func (c *Client) editTopLevelTransactionCategory(b *gotgbot.Bot, ctx *ext.Context, transaction model.Transaction) error {
	// Create keyboard with appropriate categories
	var keyboard [][]gotgbot.KeyboardButton

	if transaction.Type == model.TypeIncome {
		keyboard = [][]gotgbot.KeyboardButton{
			{{Text: "Cancel"}},
			{{Text: "Salary"}},
			{{Text: "OtherIncomes"}},
		}
	} else {
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

	// Set user state
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateTopLevelEditingTransactionCategory
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send keyboard
	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		fmt.Sprintf("Select a new category for the transaction:\n\nCurrent: <b>%s</b> - %.2f€ (%s)",
			transaction.Category,
			transaction.Amount,
			transaction.Date.Format("02-01-2006")),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
		},
	)
	if err != nil {
		return err
	}

	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Choose a category:", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard:        keyboard,
			OneTimeKeyboard: true,
			IsPersistent:    false,
			ResizeKeyboard:  true,
		},
	})

	return err
}

func (c *Client) EditTransactionCategoryConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
	requestCtx := ctx.Request.Context()
	if requestCtx == nil {
		requestCtx = context.Background()
	}
	startTime := time.Now()
	var operationErr error

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if operationErr != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("operation.type", "edit"),
			attribute.String("edit.field", "category"),
			attribute.String("status", status),
		}
		c.TransactionOperationsCounter.Add(requestCtx, 1, metric.WithAttributes(attrs...))
		c.TransactionOperationDuration.Record(requestCtx, duration, metric.WithAttributes(attrs...))
	}()

	_, u := c.getUserFromContext(ctx)
	user, operationErr := c.authAndGetUser(u)
	if operationErr != nil {
		return operationErr
	}

	// Get transaction ID from session
	transactionID, operationErr := strconv.ParseInt(user.Session.Body, 10, 64)
	if operationErr != nil {
		operationErr = fmt.Errorf("invalid transaction ID in session: %v", operationErr)
		return operationErr
	}

	// Get the transaction
	transaction, operationErr := c.Repositories.Transactions.GetByID(transactionID)
	if operationErr != nil {
		operationErr = fmt.Errorf("failed to get transaction: %w", operationErr)
		return operationErr
	}

	// Get new category from message
	newCategory := ctx.Message.Text

	// Verify it's a valid category
	if !model.IsValidTransactionCategory(newCategory) {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Invalid category. Please select a valid category.",
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
			},
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for invalid category: %v", sendErr)
		}
		operationErr = fmt.Errorf("invalid category selected: %s", newCategory)
		return operationErr
	}

	// Check if category is valid for the transaction type
	isIncome := transaction.Type == model.TypeIncome
	isIncomeCategory := newCategory == "Salary" || newCategory == "OtherIncomes"

	if isIncome != isIncomeCategory {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			fmt.Sprintf("Cannot change between expense and income categories. Please select a valid %s category.", transaction.Type),
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
			},
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for category type mismatch: %v", sendErr)
		}
		operationErr = fmt.Errorf("category type mismatch for: %s", newCategory)
		return operationErr
	}

	// Update the transaction
	oldCategory := transaction.Category
	transaction.Category = model.TransactionCategory(newCategory)

	operationErr = c.Repositories.Transactions.Update(&transaction)
	if operationErr != nil {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
			},
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for transaction update failure: %v", sendErr)
		}
		return operationErr
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	updateUserErr := c.Repositories.Users.Update(&user)
	if updateUserErr != nil {
		// Log this error, but the main operation (transaction update) was successful.
		// Consider if this should affect the 'status' for the metric. For now, it doesn't.
		c.Logger.Errorf("Failed to update user data after category edit: %v", updateUserErr)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, sendErr := b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Category updated successfully!\n\nChanged from <b>%s</b> to <b>%s</b>",
			emoji, oldCategory, transaction.Category),
		&gotgbot.SendMessageOpts{
			ParseMode:   "HTML",
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
		},
	)
	if sendErr != nil {
		c.Logger.Warnf("Error sending category update confirmation: %v", sendErr)
	}

	return operationErr // should be nil if successful
}

func (c *Client) editTopLevelTransactionDescription(b *gotgbot.Bot, ctx *ext.Context, transaction model.Transaction) error {
	// Set user state
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateTopLevelEditingTransactionDescription
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send message asking for new amount
	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		fmt.Sprintf("Enter a new description for the transaction:\n\nCurrent: <b>%s</b> (%s).",
			transaction.Description, transaction.Category),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
		},
	)

	return err
}

func (c *Client) editTopLevelTransactionAmount(b *gotgbot.Bot, ctx *ext.Context, transaction model.Transaction) error {
	// Set user state
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateTopLevelEditingTransactionAmount
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send message asking for new amount
	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		fmt.Sprintf("Enter a new amount for the transaction:\n\nCurrent: <b>%s</b> - %.2f€ (%s)",
			transaction.Category,
			transaction.Amount,
			transaction.Date.Format("02-01-2006")),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
		},
	)

	return err
}

func (c *Client) EditTransactionDescriptionConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
	requestCtx := ctx.Request.Context()
	if requestCtx == nil {
		requestCtx = context.Background()
	}
	startTime := time.Now()
	var operationErr error

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if operationErr != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("operation.type", "edit"),
			attribute.String("edit.field", "description"),
			attribute.String("status", status),
		}
		c.TransactionOperationsCounter.Add(requestCtx, 1, metric.WithAttributes(attrs...))
		c.TransactionOperationDuration.Record(requestCtx, duration, metric.WithAttributes(attrs...))
	}()

	_, u := c.getUserFromContext(ctx)
	user, operationErr := c.authAndGetUser(u)
	if operationErr != nil {
		return operationErr
	}

	// Get transaction ID from session
	transactionID, operationErr := strconv.ParseInt(user.Session.Body, 10, 64)
	if operationErr != nil {
		operationErr = fmt.Errorf("invalid transaction ID in session: %v", operationErr)
		return operationErr
	}

	// Get the transaction
	transaction, operationErr := c.Repositories.Transactions.GetByID(transactionID)
	if operationErr != nil {
		operationErr = fmt.Errorf("failed to get transaction: %w", operationErr)
		return operationErr
	}

	oldDescription := transaction.Description
	transaction.Description = strings.TrimSpace(ctx.Message.Text)
	if transaction.Description == "" {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Description cannot be empty.",
			nil,
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for empty description: %v", sendErr)
		}
		operationErr = fmt.Errorf("description cannot be empty")
		return operationErr
	}

	operationErr = c.Repositories.Transactions.Update(&transaction)
	if operationErr != nil {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			nil,
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for transaction update failure: %v", sendErr)
		}
		return operationErr
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	updateUserErr := c.Repositories.Users.Update(&user)
	if updateUserErr != nil {
		c.Logger.Errorf("Failed to update user data after description edit: %v", updateUserErr)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, sendErr := b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Description updated successfully!\n\nChanged from <b>%s</b> to <b>%s</b>",
			emoji, oldDescription, transaction.Description),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		},
	)
	if sendErr != nil {
		c.Logger.Warnf("Error sending description update confirmation: %v", sendErr)
	}

	return operationErr
}

func (c *Client) EditTransactionAmountConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
	requestCtx := ctx.Request.Context()
	if requestCtx == nil {
		requestCtx = context.Background()
	}
	startTime := time.Now()
	var operationErr error

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if operationErr != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("operation.type", "edit"),
			attribute.String("edit.field", "amount"),
			attribute.String("status", status),
		}
		c.TransactionOperationsCounter.Add(requestCtx, 1, metric.WithAttributes(attrs...))
		c.TransactionOperationDuration.Record(requestCtx, duration, metric.WithAttributes(attrs...))
	}()

	_, u := c.getUserFromContext(ctx)
	user, operationErr := c.authAndGetUser(u)
	if operationErr != nil {
		return operationErr
	}

	// Get transaction ID from session
	transactionID, operationErr := strconv.ParseInt(user.Session.Body, 10, 64)
	if operationErr != nil {
		operationErr = fmt.Errorf("invalid transaction ID in session: %v", operationErr)
		return operationErr
	}

	// Get the transaction
	transaction, operationErr := c.Repositories.Transactions.GetByID(transactionID)
	if operationErr != nil {
		operationErr = fmt.Errorf("failed to get transaction: %w", operationErr)
		return operationErr
	}

	// Parse new amount from message
	amountStr := strings.TrimSpace(ctx.Message.Text)
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	newAmount, parseErr := strconv.ParseFloat(amountStr, 64)
	if parseErr != nil {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Invalid amount. Please enter a valid number.",
			nil,
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for invalid amount: %v", sendErr)
		}
		operationErr = fmt.Errorf("invalid amount format: %v", parseErr)
		return operationErr
	}

	if newAmount <= 0 {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Amount must be greater than zero.",
			nil,
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for non-positive amount: %v", sendErr)
		}
		operationErr = fmt.Errorf("amount must be greater than zero")
		return operationErr
	}

	// Update the transaction
	oldAmount := transaction.Amount
	transaction.Amount = newAmount

	operationErr = c.Repositories.Transactions.Update(&transaction)
	if operationErr != nil {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			nil,
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for transaction update failure: %v", sendErr)
		}
		return operationErr
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	updateUserErr := c.Repositories.Users.Update(&user)
	if updateUserErr != nil {
		c.Logger.Errorf("Failed to update user data after amount edit: %v", updateUserErr)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, sendErr := b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Amount updated successfully!\n\nChanged from <b>%.2f€</b> to <b>%.2f€</b>",
			emoji, oldAmount, transaction.Amount),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		},
	)
	if sendErr != nil {
		c.Logger.Warnf("Error sending amount update confirmation: %v", sendErr)
	}

	return operationErr
}

// editTransactionDate prompts for a new date
func (c *Client) editTopLevelTransactionDate(b *gotgbot.Bot, ctx *ext.Context, transaction model.Transaction) error {
	// Set user state
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateTopLevelEditingTransactionDate
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send message asking for new date
	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		fmt.Sprintf("Enter a new date for the transaction (e.g. dd-mm-yyyy, dd/mm, etc):\n\nCurrent: <b>%s</b> - %.2f€ (%s)",
			transaction.Category,
			transaction.Amount,
			transaction.Date.Format("02-01-2006")),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
		},
	)

	return err
}

func (c *Client) EditTransactionDateConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
	requestCtx := ctx.Request.Context()
	if requestCtx == nil {
		requestCtx = context.Background()
	}
	startTime := time.Now()
	var operationErr error

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if operationErr != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("operation.type", "edit"),
			attribute.String("edit.field", "date"),
			attribute.String("status", status),
		}
		c.TransactionOperationsCounter.Add(requestCtx, 1, metric.WithAttributes(attrs...))
		c.TransactionOperationDuration.Record(requestCtx, duration, metric.WithAttributes(attrs...))
	}()

	_, u := c.getUserFromContext(ctx)
	user, operationErr := c.authAndGetUser(u)
	if operationErr != nil {
		return operationErr
	}

	// Get transaction ID from session
	transactionID, operationErr := strconv.ParseInt(user.Session.Body, 10, 64)
	if operationErr != nil {
		operationErr = fmt.Errorf("invalid transaction ID in session: %v", operationErr)
		return operationErr
	}

	// Get the transaction
	transaction, operationErr := c.Repositories.Transactions.GetByID(transactionID)
	if operationErr != nil {
		operationErr = fmt.Errorf("failed to get transaction: %w", operationErr)
		return operationErr
	}

	// Get date from DD-MM-YYYY to date
	newDate, parseErr := utils.ParseDate(ctx.Message.Text)
	if parseErr != nil {
		c.Logger.Warnf("failed to parse date: %v", parseErr)
		b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid date, please try again.", nil)
		operationErr = fmt.Errorf("failed to parse date: %v", parseErr)
		return operationErr
	}

	if newDate.After(time.Now()) {
		b.SendMessage(ctx.EffectiveSender.ChatId, "I don't support future dates, please try again.", nil)
		operationErr = fmt.Errorf("invalid date (future): %s", ctx.Message.Text)
		return operationErr
	}

	// Update the transaction
	oldDate := transaction.Date
	transaction.Date = newDate

	operationErr = c.Repositories.Transactions.Update(&transaction)
	if operationErr != nil {
		_, sendErr := b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			nil,
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending message for transaction update failure: %v", sendErr)
		}
		return operationErr
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	updateUserErr := c.Repositories.Users.Update(&user)
	if updateUserErr != nil {
		c.Logger.Errorf("Failed to update user data after date edit: %v", updateUserErr)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, sendErr := b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Date updated successfully!\n\nChanged from <b>%s</b> to <b>%s</b>",
			emoji,
			oldDate.Format("02-01-2006"),
			transaction.Date.Format("02-01-2006")),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		},
	)
	if sendErr != nil {
		c.Logger.Warnf("Error sending date update confirmation: %v", sendErr)
	}
	return operationErr
}

// showEditableTransactionPage displays a paginated list of all user transactions for editing
func (c *Client) showEditableTransactionPage(b *gotgbot.Bot, ctx *ext.Context, user model.User, offset int) error {
	limit := 5

	// Get all user transactions with pagination
	transactions, total, err := c.Repositories.Transactions.GetUserTransactionsPaginated(
		user.TgID,
		offset,
		limit,
	)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	if total == 0 {
		// No transactions found
		message := "You don't have any transactions to edit."

		if ctx.CallbackQuery != nil {
			_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{})
			return err
		} else {
			_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, nil)
			return err
		}
	}

	// Format transactions
	message := formatEditableTransactions(transactions, offset, int(total))

	// Create pagination keyboard with numbered buttons for editing
	keyboard := createEditPaginationKeyboard(transactions, offset, limit, int(total))

	// Send or update message
	if ctx.CallbackQuery != nil {
		_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	} else {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}
}

// formatEditableTransactions formats the transactions for display in the editing interface
func formatEditableTransactions(transactions []model.Transaction, offset, total int) string {
	var msg strings.Builder
	msg.WriteString("<b>✏️ Edit Transaction</b>\n")
	msg.WriteString("Select a transaction to edit:\n")
	msg.WriteString(fmt.Sprintf("Showing %d-%d of %d transactions\n\n", offset+1, offset+len(transactions), total))

	for i, t := range transactions {
		// Choose emoji based on transaction type
		emoji := "💰"
		if t.Type == model.TypeExpense {
			emoji = "💸"
		}

		msg.WriteString(fmt.Sprintf("<b>%d.</b> %s <b>%s</b> - %.2f€\n",
			i+1,
			emoji,
			t.Category,
			t.Amount,
		))
		msg.WriteString(fmt.Sprintf("   📅 %s\n", t.Date.Format("02-01-2006")))

		if t.Description != "" {
			msg.WriteString(fmt.Sprintf("   📝 %s\n", t.Description))
		}
		msg.WriteString("\n")
	}

	msg.WriteString("\nClick on a number to edit the corresponding transaction.")
	return msg.String()
}

// createEditPaginationKeyboard creates a keyboard with numbered buttons for editing transactions
func createEditPaginationKeyboard(transactions []model.Transaction, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Create number buttons for each transaction (up to 5 per row)
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		button := gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("edit.select.%d", t.ID),
		}
		row = append(row, button)

		// Create a new row after 5 buttons
		if len(row) == 5 {
			keyboard = append(keyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}

	// Add any remaining buttons
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// Navigation buttons (previous, cancel, next)
	var navigationRow []gotgbot.InlineKeyboardButton

	// Next page button (for older transactions)
	if offset+limit < total {
		nextOffset := offset + limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "⬅️ Previous",
			CallbackData: fmt.Sprintf("edit.page.%d", nextOffset),
		})
	}

	// Cancel button
	navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
		Text:         "❌ Cancel",
		CallbackData: "transactions.cancel",
	})

	// Previous page button (for newer transactions)
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "Next ➡️",
			CallbackData: fmt.Sprintf("edit.page.%d", prevOffset),
		})
	}

	keyboard = append(keyboard, navigationRow)
	return keyboard
}

// showEditOptions displays the options to edit a specific transaction
func (c *Client) showEditOptions(b *gotgbot.Bot, ctx *ext.Context, transaction model.Transaction) error {
	// Choose emoji based on transaction type
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	// Format message
	message := fmt.Sprintf("<b>✏️ Edit Transaction</b>\n\n%s <b>%s</b> - %.2f€\n📅 %s\n",
		emoji,
		transaction.Category,
		transaction.Amount,
		transaction.Date.Format("02-01-2006"),
	)

	if transaction.Description != "" {
		message += fmt.Sprintf("📝 %s\n", transaction.Description)
	}

	message += "\nSelect what you want to edit:"

	// Create keyboard with edit options
	keyboard := [][]gotgbot.InlineKeyboardButton{
		{
			{
				Text:         "🔖 Description",
				CallbackData: "edit.field.description",
			},
			{
				Text:         "✏️ Category",
				CallbackData: "edit.field.category",
			},
		},
		{
			{
				Text:         "💲 Amount",
				CallbackData: "edit.field.amount",
			},
			{
				Text:         "📅 Date",
				CallbackData: "edit.field.date",
			},
		},
		{
			{
				Text:         "❌ Cancel",
				CallbackData: "transactions.cancel",
			},
		},
	}

	// Send or update message
	_, _, err := ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	})

	return err
}
