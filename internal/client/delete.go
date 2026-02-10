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

	// Reset user state and start delete search flow
	user.Session.State = model.StateSelectingDeleteSearchCategory
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Show category selection for delete search
	return c.showDeleteSearchCategorySelection(b, ctx)
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

// ShowDeleteConfirmation handles displaying the confirmation message before deletion
func (c *Client) ShowDeleteConfirmation(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: delete.showconfirm.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	// Get the transaction details
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

	// Format transaction details for confirmation message
	emoji := utils.GetCategoryEmoji(transaction.Category)
	message := fmt.Sprintf(
		"Are you sure you want to delete this transaction?\n\n<b>%s</b> - %.2f€\n%s %s\n📅 %s",
		transaction.Description,
		transaction.Amount,
		emoji,
		transaction.Category,
		transaction.Date.Format("02-01-2006"),
	)

	// Create confirmation keyboard
	keyboard := [][]gotgbot.InlineKeyboardButton{
		{
			{
				Text:         "✅ Confirm Delete",
				CallbackData: fmt.Sprintf("delete.confirm.%d", transaction.ID),
			},
			{
				Text:         "❌ Cancel",
				CallbackData: "delete.page.0", // Go back to the first page of deletable transactions
			},
		},
	}

	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		message,
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		},
	)
	return err
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
					},
				},
			},
		},
	)

	return err
}

// showDeletableTransactionPage displays a paginated list of all user transactions
func (c *Client) showDeletableTransactionPage(b *gotgbot.Bot, ctx *ext.Context, user model.User, offset int) error {
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

	msg.WriteString("\nTap a number to delete.")
	return msg.String()
}

// createDeletionPaginationKeyboard creates a keyboard with numbered buttons for deleting transactions
func createDeletionPaginationKeyboard(transactions []model.Transaction, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Create number buttons (up to 10 per row)
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("delete.showconfirm.%d", t.ID),
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
				CallbackData: fmt.Sprintf("delete.page.%d", offset+limit),
			})
		}
		currentPage := (offset / limit) + 1
		totalPages := (total + limit - 1) / limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
			CallbackData: "delete.noop",
		})
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
	}

	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{{
		Text: "❌ Cancel", CallbackData: "transactions.cancel",
	}})
	return keyboard
}

// showDeleteSearchCategorySelection displays the category selection keyboard for delete search
func (c *Client) showDeleteSearchCategorySelection(b *gotgbot.Bot, ctx *ext.Context) error {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Add "All" option first
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔍 All Categories",
			CallbackData: "delete.search.category.all",
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
				CallbackData: fmt.Sprintf("delete.search.category.%s", cat),
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
			CallbackData: fmt.Sprintf("delete.search.category.%s", expenseCategories[i]),
		})

		// Add second button if exists
		if i+1 < len(expenseCategories) {
			emoji2 := utils.GetCategoryEmoji(expenseCategories[i+1])
			row = append(row, gotgbot.InlineKeyboardButton{
				Text:         fmt.Sprintf("%s %s", emoji2, expenseCategories[i+1]),
				CallbackData: fmt.Sprintf("delete.search.category.%s", expenseCategories[i+1]),
			})
		}

		keyboard = append(keyboard, row)
	}

	// Add cancel button
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "❌ Cancel",
			CallbackData: "delete.search.cancel",
		},
	})

	message := "🗑️ <b>Delete Transaction</b>\n\nFirst, select a category to search in:"

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

// DeleteSearchCategorySelected handles category selection for delete search
func (c *Client) DeleteSearchCategorySelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: delete.search.category.CATEGORY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	category := parts[3]

	// Store selected category in session
	user.Session.State = model.StateEnteringDeleteSearchQuery
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
		fmt.Sprintf("🗑️ Searching in <b>%s</b>\n\nEnter your search text or tap Show All:", categoryText),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "📋 Show All",
							CallbackData: "delete.search.showall",
						},
					},
					{
						{
							Text:         "❌ Cancel",
							CallbackData: "delete.search.cancel",
						},
					},
				},
			},
		},
	)

	return err
}

// DeleteSearchQueryEntered handles the search query input for delete
func (c *Client) DeleteSearchQueryEntered(b *gotgbot.Bot, ctx *ext.Context) error {
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

	// Perform search and show results for deleting
	return c.showDeleteSearchResults(b, ctx, user, category, searchQuery, 0)
}

// showDeleteSearchResults displays paginated search results for deleting
func (c *Client) showDeleteSearchResults(b *gotgbot.Bot, ctx *ext.Context, user model.User, category, searchQuery string, offset int) error {
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
								CallbackData: "delete.search.new",
							},
							{
								Text:         "🏠 Home",
								CallbackData: "delete.search.home",
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
								CallbackData: "delete.search.new",
							},
							{
								Text:         "🏠 Home",
								CallbackData: "delete.search.home",
							},
						},
					},
				},
			})
		}
		return err
	}

	// Format delete search results (similar to original delete page format)
	message := formatDeleteSearchResults(transactions, searchQuery, category, offset, int(total))

	// Create pagination keyboard with numbered buttons for deleting
	keyboard := createDeleteSearchPaginationKeyboard(transactions, category, searchQuery, offset, limit, int(total))

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

// formatDeleteSearchResults formats the search results for delete display
func formatDeleteSearchResults(transactions []model.Transaction, searchQuery, category string, offset, total int) string {
	var msg strings.Builder

	msg.WriteString("🗑️ <b>Delete Transaction</b>\n")

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

	msg.WriteString("\nTap a number to delete.")
	return msg.String()
}

// createDeleteSearchPaginationKeyboard creates pagination buttons for delete search results
func createDeleteSearchPaginationKeyboard(transactions []model.Transaction, category, searchQuery string, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Numbered selection buttons (up to 10 per row)
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("delete.search.select.%d", t.ID),
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
				CallbackData: fmt.Sprintf("delete.search.page.%s.%d.%s", category, max(offset-limit, 0), searchQuery),
			})
		}
		currentPage := (offset / limit) + 1
		totalPages := (total + limit - 1) / limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
			CallbackData: "delete.search.noop",
		})
		if offset+limit < total {
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "Next ➡️",
				CallbackData: fmt.Sprintf("delete.search.page.%s.%d.%s", category, offset+limit, searchQuery),
			})
		}
		keyboard = append(keyboard, navigationRow)
	}

	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: "🔍 New Search", CallbackData: "delete.search.new"},
		{Text: "🏠 Home", CallbackData: "delete.search.home"},
	})

	return keyboard
}

// DeleteSearchResultsPage handles pagination for delete search results
func (c *Client) DeleteSearchResultsPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: delete.search.page.CATEGORY.OFFSET.QUERY)
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

	return c.showDeleteSearchResults(b, ctx, user, category, searchQuery, offset)
}

// DeleteSearchTransactionSelected handles transaction selection from search results and shows delete confirmation
func (c *Client) DeleteSearchTransactionSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: delete.search.select.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	// Get the transaction details
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

	// Format transaction details for confirmation message
	emoji := utils.GetCategoryEmoji(transaction.Category)
	message := fmt.Sprintf(
		"Are you sure you want to delete this transaction?\n\n<b>%s</b> - %.2f€\n%s %s\n📅 %s",
		transaction.Description,
		transaction.Amount,
		emoji,
		transaction.Category,
		transaction.Date.Format("02-01-2006"),
	)

	// Create confirmation keyboard
	keyboard := [][]gotgbot.InlineKeyboardButton{
		{
			{
				Text:         "✅ Confirm Delete",
				CallbackData: fmt.Sprintf("delete.confirm.%d", transaction.ID),
			},
			{
				Text:         "❌ Cancel",
				CallbackData: "delete.search.new", // Go back to new search
			},
		},
	}

	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		message,
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		},
	)
	return err
}

// DeleteSearchCancel handles delete search cancellation
func (c *Client) DeleteSearchCancel(b *gotgbot.Bot, ctx *ext.Context) error {
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

// DeleteSearchHome returns to home screen
func (c *Client) DeleteSearchHome(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.SendHomeKeyboard(b, ctx, "What can I do for you?")
}

// DeleteSearchNew starts a new delete search
func (c *Client) DeleteSearchNew(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.DeleteTransactions(b, ctx)
}

// DeleteSearchNoop handles no-op callbacks (like page indicators)
func (c *Client) DeleteSearchNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// DeleteNoop handles no-op callbacks for the recent-transactions pagination
func (c *Client) DeleteNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// DeleteSearchShowAll handles "Show All" — searches with wildcard
func (c *Client) DeleteSearchShowAll(b *gotgbot.Bot, ctx *ext.Context) error {
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

	return c.showDeleteSearchResults(b, ctx, user, category, "%", 0)
}
