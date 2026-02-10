package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cashout/internal/model"
	"cashout/internal/utils"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// EditTransactions handles the /edit command
func (c *Client) EditTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user state and start edit search flow
	user.Session.State = model.StateSelectingEditSearchCategory
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Show category selection for edit search
	return c.showEditSearchCategorySelection(b, ctx)
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

// EditDone handles the completion of editing a transaction
func (c *Client) EditDone(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user session state
	user.Session.State = model.StateNormal
	user.Session.Body = ""

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Clear any inline keyboard
	if ctx.CallbackQuery != nil {
		_, _, err = ctx.CallbackQuery.Message.EditReplyMarkup(
			b,
			&gotgbot.EditMessageReplyMarkupOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{},
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to clear inline keyboard: %w", err)
		}
	}

	// Send confirmation message and home
	return c.SendHomeKeyboard(b, ctx, "Editing completed!")
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
			{{Text: "Pets"}},
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
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Cancel",
							CallbackData: "transactions.cancel",
						},
					},
				},
			},
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

	// Get new category from message
	newCategory := ctx.Message.Text

	// Verify it's a valid category
	if !model.IsValidTransactionCategory(newCategory) {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Invalid category. Please select a valid category.",
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
			},
		)
		return err
	}

	// Check if category is valid for the transaction type
	isIncome := transaction.Type == model.TypeIncome
	isIncomeCategory := newCategory == "Salary" || newCategory == "OtherIncomes"

	if isIncome != isIncomeCategory {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			fmt.Sprintf("Cannot change between expense and income categories. Please select a valid %s category.", transaction.Type),
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
			},
		)
		return err
	}

	// Update the transaction
	oldCategory := transaction.Category
	transaction.Category = model.TransactionCategory(newCategory)

	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
			},
		)
		return err
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, err = b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Category updated successfully!\n\nChanged from <b>%s</b> to <b>%s</b>",
			emoji, oldCategory, transaction.Category),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Keep Editing",
							CallbackData: fmt.Sprintf("edit.select.%d", transaction.ID),
						},
						{
							Text:         "Done",
							CallbackData: "edit.done",
						},
					},
				},
			},
		},
	)

	return err
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
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Cancel",
							CallbackData: "transactions.cancel",
						},
					},
				},
			},
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
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Cancel",
							CallbackData: "transactions.cancel",
						},
					},
				},
			},
		},
	)

	return err
}

func (c *Client) EditTransactionDescriptionConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
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

	oldDescription := transaction.Description
	transaction.Description = strings.TrimSpace(ctx.Message.Text)
	if transaction.Description == "" {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Description cannot be empty.",
			nil,
		)
		return err
	}

	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			nil,
		)
		return err
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, err = b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Description updated successfully!\n\nChanged from <b>%s</b> to <b>%s</b>",
			emoji, oldDescription, transaction.Description),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Keep Editing",
							CallbackData: fmt.Sprintf("edit.select.%d", transaction.ID),
						},
						{
							Text:         "Done",
							CallbackData: "edit.done",
						},
					},
				},
			},
		},
	)

	return err
}

func (c *Client) EditTransactionAmountConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
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
	oldAmount := transaction.Amount
	transaction.Amount = newAmount

	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			nil,
		)
		return err
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, err = b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Amount updated successfully!\n\nChanged from <b>%.2f€</b> to <b>%.2f€</b>",
			emoji, oldAmount, transaction.Amount),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Keep Editing",
							CallbackData: fmt.Sprintf("edit.select.%d", transaction.ID),
						},
						{
							Text:         "Done",
							CallbackData: "edit.done",
						},
					},
				},
			},
		},
	)

	return err
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
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Cancel",
							CallbackData: "transactions.cancel",
						},
					},
				},
			},
		},
	)

	return err
}

func (c *Client) EditTransactionDateConfirm(b *gotgbot.Bot, ctx *ext.Context) error {
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

	// Get date from DD-MM-YYYY to date
	newDate, err := utils.ParseDate(ctx.Message.Text)
	if err != nil {
		fmt.Printf("failed to parse date: %v\n", err)
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid date, please try again.", nil)
		return err
	}

	if newDate.After(time.Now()) {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "I don't support future dates, please try again.", nil)
		if err != nil {
			return err
		}
		return fmt.Errorf("invalid date: %s", ctx.Message.Text)
	}

	// Update the transaction
	oldDate := transaction.Date
	transaction.Date = newDate

	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Failed to update transaction. Please try again.",
			nil,
		)
		return err
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Send confirmation
	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	_, err = b.SendMessage(
		ctx.EffectiveSender.ChatId,
		fmt.Sprintf("%s Date updated successfully!\n\nChanged from <b>%s</b> to <b>%s</b>",
			emoji,
			oldDate.Format("02-01-2006"),
			transaction.Date.Format("02-01-2006")),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "Keep Editing",
							CallbackData: fmt.Sprintf("edit.select.%d", transaction.ID),
						},
						{
							Text:         "Done",
							CallbackData: "edit.done",
						},
					},
				},
			},
		},
	)

	return err
}

// showEditableTransactionPage displays a paginated list of all user transactions for editing
func (c *Client) showEditableTransactionPage(b *gotgbot.Bot, ctx *ext.Context, user model.User, offset int) error {
	limit := 10

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
	msg.WriteString(fmt.Sprintf("Showing %d–%d of %d\n\n", offset+1, offset+len(transactions), total))

	for i, t := range transactions {
		emoji := utils.GetCategoryEmoji(t.Category)
		sign := "-"
		if t.Type == model.TypeIncome {
			sign = "+"
		}
		msg.WriteString(fmt.Sprintf("%d. %s %s · %s€%.2f · %s\n",
			i+1, emoji, t.Description, sign, t.Amount, t.Date.Format("02/01")))
	}

	msg.WriteString("\nTap a number to edit.")
	return msg.String()
}

// createEditPaginationKeyboard creates a keyboard with numbered buttons for editing transactions
func createEditPaginationKeyboard(transactions []model.Transaction, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Create number buttons (up to 10 per row)
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("edit.select.%d", t.ID),
		})
		if len(row) == 5 {
			keyboard = append(keyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// Navigation row with page indicator
	if total > 0 {
		var navigationRow []gotgbot.InlineKeyboardButton
		if offset+limit < total {
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "⬅️ Previous",
				CallbackData: fmt.Sprintf("edit.page.%d", offset+limit),
			})
		}
		currentPage := (offset / limit) + 1
		totalPages := (total + limit - 1) / limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
			CallbackData: "edit.noop",
		})
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
	}

	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{{
		Text: "❌ Cancel", CallbackData: "transactions.cancel",
	}})
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

// showEditSearchCategorySelection displays the category selection keyboard for edit search
func (c *Client) showEditSearchCategorySelection(b *gotgbot.Bot, ctx *ext.Context) error {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Add "All" option first
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔍 All Categories",
			CallbackData: "edit.search.category.all",
		},
	})

	// Add income categories
	incomeCategories := []model.TransactionCategory{
		model.CategorySalary,
		model.CategoryOtherIncomes,
	}

	for _, cat := range incomeCategories {
		emoji := utils.GetCategoryEmoji(cat)
		keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s %s", emoji, cat),
				CallbackData: fmt.Sprintf("edit.search.category.%s", cat),
			},
		})
	}

	// Add expense categories in rows of 2
	expenseCategories := []model.TransactionCategory{
		model.CategoryCar,
		model.CategoryClothes,
		model.CategoryGrocery,
		model.CategoryHouse,
		model.CategoryBills,
		model.CategoryEntertainment,
		model.CategorySport,
		model.CategoryEatingOut,
		model.CategoryTransport,
		model.CategoryLearning,
		model.CategoryToiletry,
		model.CategoryHealth,
		model.CategoryTech,
		model.CategoryGifts,
		model.CategoryTravel,
		model.CategoryPets,
		model.CategoryOtherExpenses,
	}

	for i := 0; i < len(expenseCategories); i += 2 {
		row := []gotgbot.InlineKeyboardButton{}

		emoji := utils.GetCategoryEmoji(expenseCategories[i])
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s %s", emoji, expenseCategories[i]),
			CallbackData: fmt.Sprintf("edit.search.category.%s", expenseCategories[i]),
		})

		// Add second button if exists
		if i+1 < len(expenseCategories) {
			emoji2 := utils.GetCategoryEmoji(expenseCategories[i+1])
			row = append(row, gotgbot.InlineKeyboardButton{
				Text:         fmt.Sprintf("%s %s", emoji2, expenseCategories[i+1]),
				CallbackData: fmt.Sprintf("edit.search.category.%s", expenseCategories[i+1]),
			})
		}

		keyboard = append(keyboard, row)
	}

	// Add cancel button
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "❌ Cancel",
			CallbackData: "edit.search.cancel",
		},
	})

	message := "✏️ <b>Edit Transaction</b>\n\nFirst, select a category to search in:"

	// Send or update message
	if ctx.CallbackQuery != nil {
		_, _, err := ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	} else {
		_, err := b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}
}

// EditSearchCategorySelected handles category selection for edit search
func (c *Client) EditSearchCategorySelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: edit.search.category.CATEGORY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	category := parts[3]

	// Store selected category in session
	user.Session.State = model.StateEnteringEditSearchQuery
	user.Session.Body = category
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Ask for search query
	categoryText := "all categories"
	if category != "all" {
		emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
		categoryText = fmt.Sprintf("%s %s", emoji, category)
	}

	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		fmt.Sprintf("✏️ Searching in <b>%s</b>\n\nEnter your search text or tap Show All:", categoryText),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "📋 Show All",
							CallbackData: "edit.search.showall",
						},
					},
					{
						{
							Text:         "❌ Cancel",
							CallbackData: "edit.search.cancel",
						},
					},
				},
			},
		},
	)

	return err
}

// EditSearchQueryEntered handles the search query input for edit
func (c *Client) EditSearchQueryEntered(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Get search query
	searchQuery := strings.TrimSpace(ctx.Message.Text)
	if searchQuery == "" {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Search query cannot be empty. Please try again.", nil)
		return err
	}

	// Get category from session
	category := user.Session.Body

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Perform search and show results for editing
	return c.showEditSearchResults(b, ctx, user, category, searchQuery, 0)
}

// showEditSearchResults displays paginated search results for editing
func (c *Client) showEditSearchResults(b *gotgbot.Bot, ctx *ext.Context, user model.User, category, searchQuery string, offset int) error {
	limit := 10

	// Perform search
	var transactions []model.Transaction
	var total int64
	var err error

	if category == "all" {
		transactions, total, err = c.Repositories.Transactions.SearchUserTransactions(
			user.TgID,
			searchQuery,
			"",
			offset,
			limit,
		)
	} else {
		transactions, total, err = c.Repositories.Transactions.SearchUserTransactions(
			user.TgID,
			searchQuery,
			category,
			offset,
			limit,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to search transactions: %w", err)
	}

	if total == 0 {
		message := fmt.Sprintf("🔍 No transactions found matching \"%s\"", searchQuery)
		if category != "all" {
			emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
			message = fmt.Sprintf("🔍 No transactions found matching \"%s\" in %s %s", searchQuery, emoji, category)
		}

		if ctx.CallbackQuery != nil {
			_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{
							{
								Text:         "🔍 New Search",
								CallbackData: "edit.search.new",
							},
							{
								Text:         "🏠 Home",
								CallbackData: "edit.search.home",
							},
						},
					},
				},
			})
		} else {
			_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{
							{
								Text:         "🔍 New Search",
								CallbackData: "edit.search.new",
							},
							{
								Text:         "🏠 Home",
								CallbackData: "edit.search.home",
							},
						},
					},
				},
			})
		}
		return err
	}

	// Format edit search results (similar to original edit page format)
	message := formatEditSearchResults(transactions, searchQuery, category, offset, int(total))

	// Create pagination keyboard with numbered buttons for editing
	keyboard := createEditSearchPaginationKeyboard(transactions, category, searchQuery, offset, limit, int(total))

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

// formatEditSearchResults formats the search results for editing display
func formatEditSearchResults(transactions []model.Transaction, searchQuery, category string, offset, total int) string {
	var msg strings.Builder

	msg.WriteString("✏️ <b>Edit Transaction</b>\n")

	if searchQuery != "%" {
		msg.WriteString(fmt.Sprintf("Query: \"%s\"", searchQuery))
	}

	if category != "all" {
		emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
		msg.WriteString(fmt.Sprintf(" in %s %s", emoji, category))
	}

	msg.WriteString(fmt.Sprintf("\nShowing %d–%d of %d\n\n", offset+1, offset+len(transactions), total))

	for i, t := range transactions {
		emoji := utils.GetCategoryEmoji(t.Category)
		sign := "-"
		if t.Type == model.TypeIncome {
			sign = "+"
		}

		desc := t.Description
		if searchQuery != "%" {
			if idx := strings.Index(strings.ToLower(desc), strings.ToLower(searchQuery)); idx != -1 {
				desc = desc[:idx] + "<b>" + desc[idx:idx+len(searchQuery)] + "</b>" + desc[idx+len(searchQuery):]
			}
		}

		msg.WriteString(fmt.Sprintf("%d. %s %s · %s€%.2f · %s\n",
			i+1, emoji, desc, sign, t.Amount, t.Date.Format("02/01")))
	}

	msg.WriteString("\nTap a number to edit.")
	return msg.String()
}

// createEditSearchPaginationKeyboard creates pagination buttons for edit search results
func createEditSearchPaginationKeyboard(transactions []model.Transaction, category, searchQuery string, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Numbered selection buttons (up to 10 per row)
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("edit.search.select.%d", t.ID),
		})
		if len(row) == 5 {
			keyboard = append(keyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// Navigation with page indicator
	if total > 0 {
		var navigationRow []gotgbot.InlineKeyboardButton
		if offset > 0 {
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "⬅️ Previous",
				CallbackData: fmt.Sprintf("edit.search.page.%s.%d.%s", category, max(offset-limit, 0), searchQuery),
			})
		}
		currentPage := (offset / limit) + 1
		totalPages := (total + limit - 1) / limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
			CallbackData: "edit.search.noop",
		})
		if offset+limit < total {
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "Next ➡️",
				CallbackData: fmt.Sprintf("edit.search.page.%s.%d.%s", category, offset+limit, searchQuery),
			})
		}
		keyboard = append(keyboard, navigationRow)
	}

	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: "🔍 New Search", CallbackData: "edit.search.new"},
		{Text: "🏠 Home", CallbackData: "edit.search.home"},
	})

	return keyboard
}

// EditSearchResultsPage handles pagination for edit search results
func (c *Client) EditSearchResultsPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: edit.search.page.CATEGORY.OFFSET.QUERY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 6 {
		return fmt.Errorf("invalid callback data format")
	}

	category := parts[3]
	offset, err := strconv.Atoi(parts[4])
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}
	searchQuery := parts[5]

	return c.showEditSearchResults(b, ctx, user, category, searchQuery, offset)
}

// EditSearchTransactionSelected handles transaction selection from search results
func (c *Client) EditSearchTransactionSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: edit.search.select.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	// Get the transaction
	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Verify the transaction belongs to the user
	if transaction.TgID != user.TgID {
		_, _, err = ctx.CallbackQuery.Message.EditText(
			b,
			"⚠️ This transaction doesn't belong to you.",
			&gotgbot.EditMessageTextOpts{},
		)
		return err
	}

	// Store transaction ID in session for editing (same as existing edit flow)
	user.Session.State = model.StateNormal
	user.Session.Body = fmt.Sprintf("%d", transactionID)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	// Show edit options for the selected transaction
	return c.showEditOptions(b, ctx, transaction)
}

// EditSearchCancel handles edit search cancellation
func (c *Client) EditSearchCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.Cancel(b, ctx)
}

// EditSearchHome returns to home screen
func (c *Client) EditSearchHome(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.SendHomeKeyboard(b, ctx, "What can I do for you?")
}

// EditSearchNew starts a new edit search
func (c *Client) EditSearchNew(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.EditTransactions(b, ctx)
}

// EditSearchNoop handles no-op callbacks (like page indicators)
func (c *Client) EditSearchNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// EditNoop handles no-op callbacks for the recent-transactions pagination
func (c *Client) EditNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// EditSearchShowAll handles "Show All" — searches with wildcard
func (c *Client) EditSearchShowAll(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	category := user.Session.Body

	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.showEditSearchResults(b, ctx, user, category, "%", 0)
}
