package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"fmt"
	"strconv"
	"strings"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
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
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: delete.confirm.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	// Get the transaction before deleting it (for confirmation message)
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

	// Delete the transaction
	err = c.Repositories.Transactions.Delete(transactionID, user.TgID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	// Create emoji and message based on transaction type
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	// Send success message
	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		fmt.Sprintf("%s Transaction deleted successfully!\n\n%s: %s - %.2f€ (%s)",
			emoji,
			transaction.Category,
			transaction.Description,
			transaction.Amount,
			transaction.Date.Format("02-01-2006"),
		),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Delete Another",
							CallbackData: "delete.page.0",
						},
						{
							Text:         "Done",
							CallbackData: "transactions.cancel",
						},
					},
				},
			},
		},
	)

	return err
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
		emoji := utils.GetCategoryEmoji(t.Category)

		// Display with relative numbering (1-5)
		msg.WriteString(fmt.Sprintf("%d. <b>%s</b> - %.2f€\n",
			i+1, // Always 1-5 regardless of page
			t.Description,
			t.Amount,
		))

		msg.WriteString(fmt.Sprintf("   %s %s\n", emoji, t.Category))
		msg.WriteString(fmt.Sprintf("   📅 %s\n", t.Date.Format("02-01-2006")))
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
			Text:         fmt.Sprintf("%d", i+1), // Always show 1-5
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
