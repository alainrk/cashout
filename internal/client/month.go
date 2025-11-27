package client

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"cashout/internal/model"
	"cashout/internal/utils"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// MonthRecap now shows year/month selection keyboard
func (c *Client) MonthRecap(b *gotgbot.Bot, ctx *ext.Context) error {
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

	// Start with current year
	currentYear := time.Now().Year()
	return c.sendMonthRecapSelectionKeyboard(b, ctx, currentYear)
}

// MonthRecapYearNavigation handles year navigation in month selection for recap
func (c *Client) MonthRecapYearNavigation(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	_, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse year from callback data (format: monthrecap.year.YYYY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	return c.sendMonthRecapSelectionKeyboard(b, ctx, year)
}

// MonthRecapSelected displays the recap for selected month
func (c *Client) MonthRecapSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: monthrecap.month.YYYY.MM)
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

	return c.showMonthRecap(b, ctx, user, year, month)
}

// Helper function to send month selection keyboard for recap
func (c *Client) sendMonthRecapSelectionKeyboard(b *gotgbot.Bot, ctx *ext.Context, year int) error {
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
			CallbackData: fmt.Sprintf("monthrecap.month.%d.%02d", year, m),
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
	if year > MinYearAllowed {
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "⬅️ Previous Year",
			CallbackData: fmt.Sprintf("monthrecap.year.%d", year-1),
		})
	}

	if year < currentYear {
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "Next Year ➡️",
			CallbackData: fmt.Sprintf("monthrecap.year.%d", year+1),
		})
	}

	if len(navigationRow) > 0 {
		keyboard = append(keyboard, navigationRow)
	}

	// Add cancel button
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "❌ Cancel",
			CallbackData: "monthrecap.cancel",
		},
	})

	// Send the keyboard
	if ctx.CallbackQuery != nil {
		// Edit existing message
		_, _, err := ctx.CallbackQuery.Message.EditText(b, fmt.Sprintf("📊 Select a month from %d for recap:", year), &gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	} else {
		// Send new message
		_, err := b.SendMessage(ctx.EffectiveSender.ChatId, fmt.Sprintf("📊 Select a month from %d for recap:", year), &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}
}

// Helper function to show the month recap for a specific month
func (c *Client) showMonthRecap(b *gotgbot.Bot, ctx *ext.Context, user model.User, year int, month int) error {
	// Get monthly totals
	totals, err := c.Repositories.Transactions.GetMonthlyTotalsInYear(user.TgID, year)
	if err != nil {
		return err
	}

	// Get category breakdown
	categoryTotals, err := c.Repositories.Transactions.GetMonthCategorizedTotals(user.TgID, year, month)
	if err != nil {
		return err
	}

	t, ok := totals[month]

	if !ok {
		txt := fmt.Sprintf("No transactions for %s %d", time.Month(month).String(), year)
		return c.sendRecapWithNavigation(b, ctx, txt, "month", year, month)
	}

	// Format the message
	var text strings.Builder
	var monthTotal float64

	// Header with month name
	text.WriteString(fmt.Sprintf("📊 <b>%s %d Summary</b>\n\n", time.Month(month).String(), year))

	// --- EXPENSES SECTION ---
	if expenseAmount, ok := t[model.TypeExpense]; ok && expenseAmount > 0 {
		monthTotal -= expenseAmount
		text.WriteString(fmt.Sprintf("💸 <b>Expenses:</b> %.2f€\n", expenseAmount))

		// Add category breakdown for expenses
		if expenseCats, ok := categoryTotals[model.TypeExpense]; ok && len(expenseCats) > 0 {
			text.WriteString("\n<b>Expense Breakdown:</b>\n")

			// Sort categories by amount (descending)
			categories := make([]struct {
				Category model.TransactionCategory
				Amount   float64
			}, 0, len(expenseCats))

			for cat, amount := range expenseCats {
				categories = append(categories, struct {
					Category model.TransactionCategory
					Amount   float64
				}{cat, amount})
			}

			sort.Slice(categories, func(i, j int) bool {
				return categories[i].Amount > categories[j].Amount
			})

			// Display each category with emoji
			for _, entry := range categories {
				emoji := utils.GetCategoryEmoji(entry.Category)
				percentage := (entry.Amount / expenseAmount) * 100
				text.WriteString(fmt.Sprintf("  %s <b>%s:</b> %.2f€ (%.1f%%)\n",
					emoji, entry.Category, entry.Amount, percentage))
			}
			text.WriteString("\n")
		}
	}

	// --- INCOME SECTION ---
	if incomeAmount, ok := t[model.TypeIncome]; ok && incomeAmount > 0 {
		monthTotal += incomeAmount
		text.WriteString(fmt.Sprintf("💰 <b>Income:</b> %.2f€\n", incomeAmount))

		// Add category breakdown for income
		if incomeCats, ok := categoryTotals[model.TypeIncome]; ok && len(incomeCats) > 0 {
			text.WriteString("\n<b>Income Breakdown:</b>\n")

			// Sort categories by amount (descending)
			categories := make([]struct {
				Category model.TransactionCategory
				Amount   float64
			}, 0, len(incomeCats))

			for cat, amount := range incomeCats {
				categories = append(categories, struct {
					Category model.TransactionCategory
					Amount   float64
				}{cat, amount})
			}

			sort.Slice(categories, func(i, j int) bool {
				return categories[i].Amount > categories[j].Amount
			})

			// Display each category with emoji
			for _, entry := range categories {
				emoji := utils.GetCategoryEmoji(entry.Category)
				percentage := (entry.Amount / incomeAmount) * 100
				text.WriteString(fmt.Sprintf("  %s <b>%s:</b> %.2f€ (%.1f%%)\n",
					emoji, entry.Category, entry.Amount, percentage))
			}
			text.WriteString("\n")
		}
	}

	// --- TOTAL BALANCE ---
	var balanceEmoji string
	if monthTotal >= 0 {
		balanceEmoji = "✅"
	} else {
		balanceEmoji = "❌"
	}

	text.WriteString(fmt.Sprintf("\n%s <b>Month Balance:</b> %.2f€", balanceEmoji, monthTotal))

	return c.sendRecapWithNavigation(b, ctx, text.String(), "month", year, month)
}
