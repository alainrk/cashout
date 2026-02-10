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

// ListTransactions displays the category selection keyboard
func (c *Client) ListTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	_, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	return c.showListCategorySelection(b, ctx)
}

// showListCategorySelection displays the category selection keyboard (mirrors search)
func (c *Client) showListCategorySelection(b *gotgbot.Bot, ctx *ext.Context) error {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// "All" option first
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔍 All Categories",
			CallbackData: "list.cat.all",
		},
	})

	// Income categories
	incomeCategories := []model.TransactionCategory{
		model.CategorySalary,
		model.CategoryOtherIncomes,
	}

	for _, cat := range incomeCategories {
		emoji := utils.GetCategoryEmoji(cat)
		keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s %s", emoji, cat),
				CallbackData: fmt.Sprintf("list.cat.%s", cat),
			},
		})
	}

	// Expense categories in rows of 2
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
			CallbackData: fmt.Sprintf("list.cat.%s", expenseCategories[i]),
		})

		if i+1 < len(expenseCategories) {
			emoji2 := utils.GetCategoryEmoji(expenseCategories[i+1])
			row = append(row, gotgbot.InlineKeyboardButton{
				Text:         fmt.Sprintf("%s %s", emoji2, expenseCategories[i+1]),
				CallbackData: fmt.Sprintf("list.cat.%s", expenseCategories[i+1]),
			})
		}

		keyboard = append(keyboard, row)
	}

	// Cancel button
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "❌ Cancel",
			CallbackData: "list.cancel",
		},
	})

	message := "📋 <b>List Transactions</b>\n\nSelect a category to browse:"

	if ctx.CallbackQuery != nil {
		_, _, err := ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}

	_, err := b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	})
	return err
}

// ListCategorySelected handles category selection and shows month picker
func (c *Client) ListCategorySelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	_, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	category := parts[2] // "all" or a category name
	currentYear := time.Now().Year()
	return c.sendMonthSelectionKeyboard(b, ctx, currentYear, category)
}

// ListYearNavigation handles year navigation in month selection
// Callback format: list.year.YYYY.CATEGORY
func (c *Client) ListYearNavigation(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	_, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid callback data format")
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	category := parts[3]
	return c.sendMonthSelectionKeyboard(b, ctx, year, category)
}

// ListMonthTransactions displays transactions for selected month
// Callback format: list.month.YYYY.MM.CATEGORY
func (c *Client) ListMonthTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery
	parts := strings.Split(query.Data, ".")
	if len(parts) != 5 {
		return fmt.Errorf("invalid callback data format")
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	month, err := strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("invalid month: %v", err)
	}

	category := parts[4]
	return c.showTransactionPage(b, ctx, user, year, month, 0, category)
}

// ListTransactionPage handles transaction pagination
// Callback format: list.page.YYYY.MM.OFFSET.CATEGORY
func (c *Client) ListTransactionPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery
	parts := strings.Split(query.Data, ".")
	if len(parts) != 6 {
		return fmt.Errorf("invalid callback data format")
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	month, err := strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("invalid month: %v", err)
	}

	offset, err := strconv.Atoi(parts[4])
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}

	category := parts[5]
	return c.showTransactionPage(b, ctx, user, year, month, offset, category)
}

// sendMonthSelectionKeyboard renders the month picker with category threaded through
func (c *Client) sendMonthSelectionKeyboard(b *gotgbot.Bot, ctx *ext.Context, year int, category string) error {
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()

	var keyboard [][]gotgbot.InlineKeyboardButton

	// Month buttons (3 per row)
	var row []gotgbot.InlineKeyboardButton
	for m := 1; m <= 12; m++ {
		monthName := time.Month(m).String()[:3]

		if year == currentYear && m > int(currentMonth) {
			continue
		}

		button := gotgbot.InlineKeyboardButton{
			Text:         monthName,
			CallbackData: fmt.Sprintf("list.month.%d.%02d.%s", year, m, category),
		}
		row = append(row, button)

		if len(row) == 3 {
			keyboard = append(keyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// Year navigation
	navigationRow := []gotgbot.InlineKeyboardButton{}
	if year > 2020 {
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "⬅️ Previous Year",
			CallbackData: fmt.Sprintf("list.year.%d.%s", year-1, category),
		})
	}

	if year < currentYear {
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "Next Year ➡️",
			CallbackData: fmt.Sprintf("list.year.%d.%s", year+1, category),
		})
	}

	if len(navigationRow) > 0 {
		keyboard = append(keyboard, navigationRow)
	}

	// Back to categories + Cancel
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔙 Back to Categories",
			CallbackData: "list.backtocategories",
		},
		{
			Text:         "Cancel",
			CallbackData: "list.cancel",
		},
	})

	// Header text
	headerCategory := "all categories"
	if category != "all" {
		emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
		headerCategory = fmt.Sprintf("%s %s", emoji, category)
	}

	text := fmt.Sprintf("📋 Browsing <b>%s</b>\nSelect a month from %d:", headerCategory, year)

	if ctx.CallbackQuery != nil {
		_, _, err := ctx.CallbackQuery.Message.EditText(b, text, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}

	_, err := b.SendMessage(ctx.EffectiveSender.ChatId, text, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	})
	return err
}

// showTransactionPage renders the paginated transaction list
func (c *Client) showTransactionPage(b *gotgbot.Bot, ctx *ext.Context, user model.User, year, month, offset int, category string) error {
	limit := 20

	// Convert "all" to empty string for DB query
	dbCategory := ""
	if category != "all" {
		dbCategory = category
	}

	transactions, total, err := c.Repositories.Transactions.GetUserTransactionsByMonthPaginated(user.TgID, year, month, offset, limit, dbCategory)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	message := formatTransactions(year, month, transactions, offset, int(total), category)
	keyboard := createPaginationKeyboard(year, month, offset, limit, int(total), category)

	if ctx.CallbackQuery != nil {
		_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}

	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	})
	return err
}

// formatTransactions formats the compact transaction list
func formatTransactions(year, month int, transactions []model.Transaction, offset, total int, category string) string {
	if len(transactions) == 0 {
		msg := fmt.Sprintf("No transactions found for %s %d", time.Month(month).String(), year)
		if category != "all" {
			emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
			msg += fmt.Sprintf(" in %s %s", emoji, category)
		}
		return msg
	}

	var msg strings.Builder

	// Header with optional category
	headerCategory := ""
	if category != "all" {
		emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
		headerCategory = fmt.Sprintf(" · %s %s", emoji, category)
	}

	msg.WriteString(fmt.Sprintf("📊 <b>%s %d</b>%s\n", time.Month(month).String(), year, headerCategory))
	msg.WriteString(fmt.Sprintf("Showing %d–%d of %d\n\n", offset+1, offset+len(transactions), total))

	for _, t := range transactions {
		emoji := utils.GetCategoryEmoji(t.Category)
		sign := "-"
		if t.Type == model.TypeIncome {
			sign = "+"
		}
		msg.WriteString(fmt.Sprintf("%s <b>%s</b> · %s€%.2f · %s\n",
			emoji, t.Description, sign, t.Amount, t.Date.Format("02/01")))
	}

	return msg.String()
}

// createPaginationKeyboard creates pagination buttons with category threaded through
func createPaginationKeyboard(year, month, offset, limit, total int, category string) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton
	var navigationRow []gotgbot.InlineKeyboardButton

	// Previous page button
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "⬅️ Previous",
			CallbackData: fmt.Sprintf("list.page.%d.%02d.%d.%s", year, month, prevOffset, category),
		})
	}

	// Page indicator
	currentPage := (offset / limit) + 1
	totalPages := (total + limit - 1) / limit
	navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
		CallbackData: "list.noop",
	})

	// Next page button
	if offset+limit < total {
		nextOffset := offset + limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "Next ➡️",
			CallbackData: fmt.Sprintf("list.page.%d.%02d.%d.%s", year, month, nextOffset, category),
		})
	}

	if len(navigationRow) > 0 {
		keyboard = append(keyboard, navigationRow)
	}

	// Back to month selection
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔙 Back to Months",
			CallbackData: fmt.Sprintf("list.year.%d.%s", year, category),
		},
	})

	return keyboard
}

// ListNoop handles no-op callbacks (like page indicators)
func (c *Client) ListNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// ListBackToCategories returns to category selection
func (c *Client) ListBackToCategories(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.showListCategorySelection(b, ctx)
}
