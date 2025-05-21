package client

import (
	"cashout/internal/model"
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

// DeleteTransactions handles the /delete command
func (c *Client) DeleteTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
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

	// Show first page of deletable transactions
	return c.showDeletableTransactionPage(b, ctx, user, 0)
}

// DeleteTransactionPage handles pagination in the transaction deletion interface
func (c *Client) DeleteTransactionPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: delete.page.OFFSET)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	offset, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}

	return c.showDeletableTransactionPage(b, ctx, user, offset)
}

// DeleteTransactionConfirm handles the confirmation callback for deleting a transaction
func (c *Client) DeleteTransactionConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
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
			attribute.String("operation.type", "delete"),
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

	query := ctx.CallbackQuery

	// Parse callback data (format: delete.confirm.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		operationErr = fmt.Errorf("invalid callback data format for delete.confirm")
		return operationErr
	}

	transactionID, parseErr := strconv.ParseInt(parts[2], 10, 64)
	if parseErr != nil {
		operationErr = fmt.Errorf("invalid transaction ID: %v", parseErr)
		return operationErr
	}

	// Get the transaction before deleting it (for confirmation message)
	// This is good practice, but if it fails, the delete operation itself hasn't failed yet.
	transaction, getErr := c.Repositories.Transactions.GetByID(transactionID)
	if getErr != nil {
		// Log this, but proceed to attempt deletion if user is authorized.
		// The delete operation itself will be the primary point of failure assessment for the metric.
		c.Logger.Warnf("Failed to get transaction %d before deletion (for confirmation message): %v", transactionID, getErr)
		// We might not have transaction.Type for the emoji later, handle this.
	}

	// Verify ownership - this is crucial. If GetByID failed, transaction might be nil.
	if transaction != nil && transaction.TgID != user.TgID {
		_, _, sendErr := ctx.CallbackQuery.Message.EditText(
			b,
			"⚠️ This transaction doesn't belong to you.",
			&gotgbot.EditMessageTextOpts{},
		)
		if sendErr != nil {
			c.Logger.Warnf("Error sending ownership error message: %v", sendErr)
		}
		operationErr = fmt.Errorf("user %d attempted to delete transaction %d owned by %d", user.TgID, transactionID, transaction.TgID)
		return operationErr
	} else if transaction == nil && getErr != nil {
		// If we couldn't get the transaction, we can't verify ownership this way.
		// The .Delete method should internally handle ownership or TgID matching for safety.
		// For now, we proceed, relying on the Delete method's own checks.
		c.Logger.Warnf("Could not verify ownership for transaction %d due to GetByID error. Proceeding with delete attempt.", transactionID)
	}


	// Delete the transaction
	operationErr = c.Repositories.Transactions.Delete(transactionID, user.TgID)
	if operationErr != nil {
		// Don't return yet, let defer handle metrics.
		// Error message will be sent by caller or not at all if this is a background task.
		// For now, we assume this function sends its own errors if needed.
		// SendMessage(ctx, b, "Failed to delete transaction.", nil) // Example if needed
		operationErr = fmt.Errorf("failed to delete transaction %d: %w", transactionID, operationErr)
		return operationErr // Return after setting operationErr so defer captures it.
	}

	// Create emoji and message based on transaction type
	emoji := "🗑" // General delete emoji
	var confirmationMessage string
	if transaction != nil { // If we successfully got the transaction details
		if transaction.Type == model.TypeExpense {
			emoji = "💸🗑"
		} else {
			emoji = "💰🗑"
		}
		confirmationMessage = fmt.Sprintf("%s Transaction deleted successfully!\n\nDetails: %s - %.2f€ on %s",
			emoji,
			transaction.Category,
			// transaction.Description, // Description can be long, keep it concise
			transaction.Amount,
			transaction.Date.Format("02-01-2006"),
		)
	} else { // Fallback message if transaction details weren't available
		confirmationMessage = fmt.Sprintf("%s Transaction (ID: %d) deleted successfully!", emoji, transactionID)
	}


	// Send success message
	_, _, sendEditErr := ctx.CallbackQuery.Message.EditText(
		b,
		confirmationMessage,
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Delete Another",
							CallbackData: "delete.page.0", // Navigate to first page of delete options
						},
						{
							Text:         "Done",
							CallbackData: "transactions.cancel", // Go to home/cancel
						},
					},
				},
			},
		},
	)
	if sendEditErr != nil {
		c.Logger.Warnf("Error sending delete confirmation message: %v", sendEditErr)
		// The delete operation was successful, but sending confirmation failed.
		// Not changing operationErr as the core DB operation succeeded.
	}

	return operationErr // should be nil if successful
}

// showDeletableTransactionPage displays a paginated list of all user transactions
func (c *Client) showDeletableTransactionPage(b *gotgbot.Bot, ctx *ext.Context, user model.User, offset int) error {
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
		message := "You don't have any transactions to delete."

		if ctx.CallbackQuery != nil {
			_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{})
			return err
		} else {
			_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, nil)
			return err
		}
	}

	// Format transactions
	message := formatDeletableTransactions(transactions, offset, int(total))

	// Create pagination keyboard with numbered buttons for deletion
	keyboard := createDeletionPaginationKeyboard(transactions, offset, limit, int(total))

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

// formatDeletableTransactions formats the transactions for display in the deletion interface
func formatDeletableTransactions(transactions []model.Transaction, offset, total int) string {
	var msg strings.Builder
	msg.WriteString("<b>🗑 Delete Transaction</b>\n")
	msg.WriteString("Select a transaction to delete:\n")
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

	msg.WriteString("\nClick on a number to delete the corresponding transaction.")
	return msg.String()
}

// createDeletionPaginationKeyboard creates a keyboard with numbered buttons for deleting transactions
func createDeletionPaginationKeyboard(transactions []model.Transaction, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Create number buttons for each transaction (up to 5 per row)
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		button := gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("delete.confirm.%d", t.ID),
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
			CallbackData: fmt.Sprintf("delete.page.%d", nextOffset),
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
			CallbackData: fmt.Sprintf("delete.page.%d", prevOffset),
		})
	}

	keyboard = append(keyboard, navigationRow)
	return keyboard
}
