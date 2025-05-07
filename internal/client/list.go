package client

import (
	"cashout/internal/model"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// ListTransactions displays the year/month selection keyboard
func (c *Client) ListTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	_, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Start with current year
	currentYear := time.Now().Year()
	return c.sendMonthSelectionKeyboard(b, ctx, currentYear)
}

// ListYearNavigation handles year navigation in month selection
func (c *Client) ListYearNavigation(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	_, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse year from callback data (format: list.year.YYYY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	return c.sendMonthSelectionKeyboard(b, ctx, year)
}

// ListMonthTransactions displays transactions for selected month
func (c *Client) ListMonthTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: list.month.YYYY.MM)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 4 {
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

	return c.showTransactionPage(b, ctx, user, year, month, 0)
}

// ListTransactionPage handles transaction pagination
func (c *Client) ListTransactionPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: list.page.YYYY.MM.OFFSET)
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

	offset, err := strconv.Atoi(parts[4])
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}

	return c.showTransactionPage(b, ctx, user, year, month, offset)
}

// Helper function to send month selection keyboard
func (c *Client) sendMonthSelectionKeyboard(b *gotgbot.Bot, ctx *ext.Context, year int) error {
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()

	var keyboard [][]gotgbot.InlineKeyboardButton

	// Create month buttons (3 months per row)
	var row []gotgbot.InlineKeyboardButton
	for m := 1; m <= 12; m++ {
		monthName := time.Month(m).String()[:3] // Short month name

		// Disable future months for current year
		if year == currentYear && m > int(currentMonth) {
			continue
		}

		button := gotgbot.InlineKeyboardButton{
			Text:         monthName,
			CallbackData: fmt.Sprintf("list.month.%d.%02d", year, m),
		}
		row = append(row, button)

		if len(row) == 3 {
			keyboard = append(keyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}

	// Add remaining months if any
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// Add navigation buttons
	navigationRow := []gotgbot.InlineKeyboardButton{}
	if year > 2020 { // Arbitrary minimum year
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "⬅️ Previous Year",
			CallbackData: fmt.Sprintf("list.year.%d", year-1),
		})
	}

	if year < currentYear {
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "Next Year ➡️",
			CallbackData: fmt.Sprintf("list.year.%d", year+1),
		})
	}

	if len(navigationRow) > 0 {
		keyboard = append(keyboard, navigationRow)
	}

	// In any case, add the cancel
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "Cancel",
			CallbackData: "list.cancel",
		},
	})

	// Send the keyboard
	if ctx.CallbackQuery != nil {
		// Edit existing message
		_, _, err := ctx.CallbackQuery.Message.EditText(b, fmt.Sprintf("Select a month from %d:", year), &gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	} else {
		// Send new message
		_, err := b.SendMessage(ctx.EffectiveSender.ChatId, fmt.Sprintf("Select a month from %d:", year), &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}
}

// Helper function to show transactions page
func (c *Client) showTransactionPage(b *gotgbot.Bot, ctx *ext.Context, user model.User, year, month, offset int) error {
	limit := 10

	transactions, total, err := c.Repositories.Transactions.GetUserTransactionsByMonthPaginated(user.TgID, year, month, offset, limit)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	// Format transactions
	message := formatTransactions(year, month, transactions, offset, int(total))

	// Create pagination keyboard
	keyboard := createPaginationKeyboard(year, month, offset, limit, int(total))

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

// Helper function to format transactions
func formatTransactions(year, month int, transactions []model.Transaction, offset, total int) string {
	if len(transactions) == 0 {
		return fmt.Sprintf("No transactions found for %s %d", time.Month(month).String(), year)
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("📊 <b>%s %d</b>\n", time.Month(month).String(), year))
	msg.WriteString(fmt.Sprintf("Showing %d-%d of %d transactions\n\n", offset+1, offset+len(transactions), total))

	for i, t := range transactions {
		// Choose emoji based on transaction type
		emoji := "💰"
		if t.Type == model.TypeExpense {
			emoji = "💸"
		}

		msg.WriteString(fmt.Sprintf("%d. %s <b>%s</b> - %.2f€\n",
			offset+i+1,
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

	return msg.String()
}

// Helper function to create pagination keyboard
func createPaginationKeyboard(year, month, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
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
			CallbackData: fmt.Sprintf("list.page.%d.%02d.%d", year, month, prevOffset),
		})
	}

	// Next page button
	if offset+limit < total {
		nextOffset := offset + limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "Next ➡️",
			CallbackData: fmt.Sprintf("list.page.%d.%02d.%d", year, month, nextOffset),
		})
	}

	if len(navigationRow) > 0 {
		keyboard = append(keyboard, navigationRow)
	}

	// Back to month selection button
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔙 Back to Months",
			CallbackData: fmt.Sprintf("list.year.%d", year),
		},
	})

	return keyboard
}
