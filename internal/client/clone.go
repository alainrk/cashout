package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// CloneTransactions handles the /clone command and home.clone callback
func (c *Client) CloneTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateSelectingCloneTransaction
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.showRecentExpensesForClone(b, ctx, user, 0)
}

// showRecentExpensesForClone displays the 10 most recent expenses for cloning
func (c *Client) showRecentExpensesForClone(b *gotgbot.Bot, ctx *ext.Context, user model.User, offset int) error {
	limit := 10

	transactions, total, err := c.Repositories.Transactions.GetUserTransactionsByTypePaginated(
		user.TgID,
		model.TypeExpense,
		offset,
		limit,
	)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	if total == 0 {
		message := "You don't have any transactions to clone. Add your first transaction!"
		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "💰 Add Income", CallbackData: "transactions.new.income"},
				{Text: "💸 Add Expense", CallbackData: "transactions.new.expense"},
			},
			{
				{Text: "🏠 Home", CallbackData: "clone.search.home"},
			},
		}
		return SendMessage(ctx, b, message, keyboard)
	}

	// Format transactions and build keyboard
	msg := formatCloneRecentExpenses(transactions, offset, int(total))
	keyboard := createCloneRecentKeyboard(transactions, offset, limit, int(total))

	return SendMessage(ctx, b, msg, keyboard)
}

// CloneTransactionSelected handles selection of a transaction to clone from the recent list
func (c *Client) CloneTransactionSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: clone.select.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	return c.cloneAndSaveTransaction(b, ctx, user, transactionID)
}

// cloneAndSaveTransaction fetches a transaction, clones it with today's date, saves it, and shows the edit UI
func (c *Client) cloneAndSaveTransaction(b *gotgbot.Bot, ctx *ext.Context, user model.User, sourceID int64) error {
	// Get the source transaction
	source, err := c.Repositories.Transactions.GetByID(sourceID)
	if err != nil {
		keyboard := [][]gotgbot.InlineKeyboardButton{
			{{Text: "🏠 Home", CallbackData: "clone.search.home"}},
		}
		return SendMessage(ctx, b, "⚠️ Transaction no longer exists.", keyboard)
	}

	// Verify ownership
	if source.TgID != user.TgID {
		keyboard := [][]gotgbot.InlineKeyboardButton{
			{{Text: "🏠 Home", CallbackData: "clone.search.home"}},
		}
		return SendMessage(ctx, b, "⚠️ This transaction doesn't belong to you.", keyboard)
	}

	// Create clone with today's date
	clone := model.Transaction{
		TgID:        user.TgID,
		Date:        time.Now(),
		Type:        source.Type,
		Category:    source.Category,
		Amount:      source.Amount,
		Currency:    source.Currency,
		Description: source.Description,
	}

	err = c.Repositories.Transactions.Add(&clone)
	if err != nil {
		return fmt.Errorf("failed to save cloned transaction: %w", err)
	}

	// Store the clone ID in session for potential edits (same as add flow)
	user.Session.State = model.StateEditingNewTransaction
	user.Session.Body = strconv.FormatInt(clone.ID, 10)
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	emoji := "💰"
	if clone.Type == model.TypeExpense {
		emoji = "💸"
	}

	msg := fmt.Sprintf("%s <b>Transaction cloned!</b>\n\n%s (€ %.2f), %s on %s",
		emoji, clone.Category, clone.Amount, clone.Description, clone.Date.Format("02-01-2006"))

	keyboard := [][]gotgbot.InlineKeyboardButton{
		{{Text: "Edit description", CallbackData: "transactions.edit.description"}},
		{{Text: "Edit category", CallbackData: "transactions.edit.category"}},
		{{Text: "Edit date", CallbackData: "transactions.edit.date"}},
		{{Text: "Edit amount", CallbackData: "transactions.edit.amount"}},
		{
			{Text: "Delete", CallbackData: fmt.Sprintf("transactions.delete.%d", clone.ID)},
			{Text: "Home", CallbackData: "transactions.home"},
		},
	}

	return SendMessage(ctx, b, msg, keyboard)
}

// CloneTransactionPage handles pagination for the recent expenses list
func (c *Client) CloneTransactionPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	offset, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}

	return c.showRecentExpensesForClone(b, ctx, user, offset)
}

// CloneNoop handles no-op callbacks (page indicator)
func (c *Client) CloneNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// --- Wizard Search (US2) ---

// CloneSearchMore shows the type selection screen for the wizard
func (c *Client) CloneSearchMore(b *gotgbot.Bot, ctx *ext.Context) error {
	keyboard := [][]gotgbot.InlineKeyboardButton{
		{
			{Text: "💸 Expenses", CallbackData: "clone.search.type.expense"},
			{Text: "💰 Incomes", CallbackData: "clone.search.type.income"},
		},
		{
			{Text: "All", CallbackData: "clone.search.type.all"},
		},
		{
			{Text: "❌ Cancel", CallbackData: "clone.search.cancel"},
		},
	}

	return SendMessage(ctx, b, "📋 <b>Clone Transaction</b>\n\nFilter by type:", keyboard)
}

// CloneSearchTypeSelected handles type selection in the wizard
func (c *Client) CloneSearchTypeSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: clone.search.type.TYPE)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	selectedType := parts[3] // "expense", "income", or "all"

	user.Session.State = model.StateSelectingCloneSearchCategory
	user.Session.Body = selectedType
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.showCloneSearchCategorySelection(b, ctx, selectedType)
}

// showCloneSearchCategorySelection displays category selection filtered by type
func (c *Client) showCloneSearchCategorySelection(b *gotgbot.Bot, ctx *ext.Context, selectedType string) error {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// "All Categories" option
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: "🔍 All Categories", CallbackData: "clone.search.category.all"},
	})

	// Income categories
	if selectedType == "income" || selectedType == "all" {
		incomeCategories := []model.TransactionCategory{
			model.CategorySalary,
			model.CategoryOtherIncomes,
		}
		for _, cat := range incomeCategories {
			emoji := utils.GetCategoryEmoji(cat)
			keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
				{Text: fmt.Sprintf("%s %s", emoji, cat), CallbackData: fmt.Sprintf("clone.search.category.%s", cat)},
			})
		}
	}

	// Expense categories
	if selectedType == "expense" || selectedType == "all" {
		expenseCategories := []model.TransactionCategory{
			model.CategoryCar, model.CategoryClothes, model.CategoryGrocery,
			model.CategoryHouse, model.CategoryBills, model.CategoryEntertainment,
			model.CategorySport, model.CategoryEatingOut, model.CategoryTransport,
			model.CategoryLearning, model.CategoryToiletry, model.CategoryHealth,
			model.CategoryTech, model.CategoryGifts, model.CategoryTravel,
			model.CategoryPets, model.CategoryOtherExpenses,
		}
		for i := 0; i < len(expenseCategories); i += 2 {
			var row []gotgbot.InlineKeyboardButton
			emoji := utils.GetCategoryEmoji(expenseCategories[i])
			row = append(row, gotgbot.InlineKeyboardButton{
				Text:         fmt.Sprintf("%s %s", emoji, expenseCategories[i]),
				CallbackData: fmt.Sprintf("clone.search.category.%s", expenseCategories[i]),
			})
			if i+1 < len(expenseCategories) {
				emoji2 := utils.GetCategoryEmoji(expenseCategories[i+1])
				row = append(row, gotgbot.InlineKeyboardButton{
					Text:         fmt.Sprintf("%s %s", emoji2, expenseCategories[i+1]),
					CallbackData: fmt.Sprintf("clone.search.category.%s", expenseCategories[i+1]),
				})
			}
			keyboard = append(keyboard, row)
		}
	}

	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: "❌ Cancel", CallbackData: "clone.search.cancel"},
	})

	return SendMessage(ctx, b, "📋 <b>Clone Transaction</b>\n\nSelect a category to search in:", keyboard)
}

// CloneSearchCategorySelected handles category selection in the wizard
func (c *Client) CloneSearchCategorySelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: clone.search.category.CATEGORY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	category := parts[3]

	user.Session.State = model.StateEnteringCloneSearchQuery
	user.Session.Body = category
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	categoryText := "all categories"
	if category != "all" {
		emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
		categoryText = fmt.Sprintf("%s %s", emoji, category)
	}

	keyboard := [][]gotgbot.InlineKeyboardButton{
		{{Text: "📋 Show All", CallbackData: "clone.search.showall"}},
		{{Text: "❌ Cancel", CallbackData: "clone.search.cancel"}},
	}

	return SendMessage(ctx, b, fmt.Sprintf("📋 Searching in <b>%s</b>\n\nEnter your search text or tap Show All:", categoryText), keyboard)
}

// CloneSearchQueryEntered handles the free-text search query input
func (c *Client) CloneSearchQueryEntered(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	searchQuery := strings.TrimSpace(ctx.Message.Text)
	if searchQuery == "" {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Search query cannot be empty. Please try again.", nil)
		return err
	}

	category := user.Session.Body

	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.showCloneSearchResults(b, ctx, user, category, searchQuery, 0)
}

// showCloneSearchResults displays paginated search results for cloning
func (c *Client) showCloneSearchResults(b *gotgbot.Bot, ctx *ext.Context, user model.User, category, searchQuery string, offset int) error {
	limit := 10

	var transactions []model.Transaction
	var total int64
	var err error

	if category == "all" {
		transactions, total, err = c.Repositories.Transactions.SearchUserTransactions(user.TgID, searchQuery, "", offset, limit)
	} else {
		transactions, total, err = c.Repositories.Transactions.SearchUserTransactions(user.TgID, searchQuery, category, offset, limit)
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
		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "🔍 New Search", CallbackData: "clone.search.new"},
				{Text: "🏠 Home", CallbackData: "clone.search.home"},
			},
		}
		return SendMessage(ctx, b, message, keyboard)
	}

	message := formatCloneSearchResults(transactions, searchQuery, category, offset, int(total))
	keyboard := createCloneSearchKeyboard(transactions, category, searchQuery, offset, limit, int(total))

	return SendMessage(ctx, b, message, keyboard)
}

// CloneSearchTransactionSelected handles transaction selection from search results
func (c *Client) CloneSearchTransactionSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: clone.search.select.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	return c.cloneAndSaveTransaction(b, ctx, user, transactionID)
}

// CloneSearchResultsPage handles pagination for search results
func (c *Client) CloneSearchResultsPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: clone.search.page.CATEGORY.OFFSET.QUERY)
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

	return c.showCloneSearchResults(b, ctx, user, category, searchQuery, offset)
}

// CloneSearchShowAll handles "Show All" — searches with wildcard
func (c *Client) CloneSearchShowAll(b *gotgbot.Bot, ctx *ext.Context) error {
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

	return c.showCloneSearchResults(b, ctx, user, category, "%", 0)
}

// CloneSearchNoop handles no-op callbacks for search pagination
func (c *Client) CloneSearchNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// CloneSearchNew starts a new clone search (restarts the flow)
func (c *Client) CloneSearchNew(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.CloneTransactions(b, ctx)
}

// CloneSearchCancel cancels the clone search and resets state
func (c *Client) CloneSearchCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.Cancel(b, ctx)
}

// CloneSearchHome returns to the home screen
func (c *Client) CloneSearchHome(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.SendHomeKeyboard(b, ctx, "What can I do for you?")
}

// --- Extracted pure functions for formatting and keyboard building (testable) ---

// formatCloneRecentExpenses formats the recent expenses list for the clone UI
func formatCloneRecentExpenses(transactions []model.Transaction, offset, total int) string {
	var msg strings.Builder
	msg.WriteString("📋 <b>Clone Transaction</b>\n")
	msg.WriteString(fmt.Sprintf("Recent expenses — %d–%d of %d\n\n", offset+1, offset+len(transactions), total))

	for i, t := range transactions {
		emoji := utils.GetCategoryEmoji(t.Category)
		msg.WriteString(fmt.Sprintf("%d. %s %s · €%.2f · %s\n",
			i+1, emoji, t.Description, t.Amount, t.Date.Format("02/01/2006")))
	}

	msg.WriteString("\nTap a number to clone it with today's date.")
	return msg.String()
}

// createCloneRecentKeyboard creates the keyboard for the recent expenses clone list
func createCloneRecentKeyboard(transactions []model.Transaction, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Numbered selection buttons (rows of 5)
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("clone.select.%d", t.ID),
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
	if total > limit {
		var navigationRow []gotgbot.InlineKeyboardButton
		if offset+limit < total {
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "⬅️ Previous",
				CallbackData: fmt.Sprintf("clone.page.%d", offset+limit),
			})
		}
		currentPage := (offset / limit) + 1
		totalPages := (total + limit - 1) / limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
			CallbackData: "clone.noop",
		})
		if offset > 0 {
			prevOffset := offset - limit
			if prevOffset < 0 {
				prevOffset = 0
			}
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "Next ➡️",
				CallbackData: fmt.Sprintf("clone.page.%d", prevOffset),
			})
		}
		keyboard = append(keyboard, navigationRow)
	}

	// Search More + Cancel
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: "🔍 Search More", CallbackData: "clone.searchmore"},
		{Text: "❌ Cancel", CallbackData: "transactions.cancel"},
	})

	return keyboard
}

// formatCloneSearchResults formats the search results for the clone UI
func formatCloneSearchResults(transactions []model.Transaction, searchQuery, category string, offset, total int) string {
	var msg strings.Builder
	msg.WriteString("📋 <b>Clone Transaction</b>\n")
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
			i+1, emoji, desc, sign, t.Amount, t.Date.Format("02/01/2006")))
	}

	msg.WriteString("\nTap a number to clone it with today's date.")
	return msg.String()
}

// createCloneSearchKeyboard creates the keyboard for the clone search results
func createCloneSearchKeyboard(transactions []model.Transaction, category, searchQuery string, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Numbered selection buttons
	var row []gotgbot.InlineKeyboardButton
	for i, t := range transactions {
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i+1),
			CallbackData: fmt.Sprintf("clone.search.select.%d", t.ID),
		})
		if len(row) == 5 {
			keyboard = append(keyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// Navigation
	if total > limit {
		var navigationRow []gotgbot.InlineKeyboardButton
		if offset > 0 {
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "⬅️ Previous",
				CallbackData: fmt.Sprintf("clone.search.page.%s.%d.%s", category, max(offset-limit, 0), searchQuery),
			})
		}
		currentPage := (offset / limit) + 1
		totalPages := (total + limit - 1) / limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
			CallbackData: "clone.search.noop",
		})
		if offset+limit < total {
			navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
				Text:         "Next ➡️",
				CallbackData: fmt.Sprintf("clone.search.page.%s.%d.%s", category, offset+limit, searchQuery),
			})
		}
		keyboard = append(keyboard, navigationRow)
	}

	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: "🔍 New Search", CallbackData: "clone.search.new"},
		{Text: "🏠 Home", CallbackData: "clone.search.home"},
	})

	return keyboard
}
