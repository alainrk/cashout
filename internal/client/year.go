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

// YearRecap now shows year selection keyboard
func (c *Client) YearRecap(b *gotgbot.Bot, ctx *ext.Context) error {
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

	// Show year selection keyboard
	return c.sendYearRecapSelectionKeyboard(b, ctx)
}

// YearRecapSelected displays the recap for selected year
func (c *Client) YearRecapSelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: yearrecap.year.YYYY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid year: %v", err)
	}

	return c.showYearRecap(b, ctx, user, year)
}

// Helper function to send year selection keyboard
func (c *Client) sendYearRecapSelectionKeyboard(b *gotgbot.Bot, ctx *ext.Context) error {
	currentYear := time.Now().Year()

	var keyboard [][]gotgbot.InlineKeyboardButton

	// Create year buttons (4 years per row)
	var row []gotgbot.InlineKeyboardButton

	// Show years from current year down to MIN_YEAR_ALLOWED
	for year := MinYearAllowed; year <= currentYear; year++ {
		button := gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", year),
			CallbackData: fmt.Sprintf("yearrecap.year.%d", year),
		}
		row = append(row, button)

		if len(row) == 4 {
			keyboard = append(keyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}

	// Add remaining years if any
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}

	// Add cancel button
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "❌ Cancel",
			CallbackData: "yearrecap.cancel",
		},
	})

	// Send the keyboard
	if ctx.CallbackQuery != nil {
		// Edit existing message
		_, _, err := ctx.CallbackQuery.Message.EditText(b, "📊 Select a year for recap:", &gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	} else {
		// Send new message
		_, err := b.SendMessage(ctx.EffectiveSender.ChatId, "📊 Select a year for recap:", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}
}

// Helper function to show the year recap for a specific year
func (c *Client) showYearRecap(b *gotgbot.Bot, ctx *ext.Context, user model.User, year int) error {
	// Get monthly totals for all months
	res, err := c.Repositories.Transactions.GetMonthlyTotalsInYear(user.TgID, year)
	if err != nil {
		return err
	}

	// Get category breakdown for the entire year
	categoryTotals, err := c.Repositories.Transactions.GetYearCategorizedTotals(user.TgID, year)
	if err != nil {
		return err
	}

	var msg strings.Builder
	var yearTotal float64
	var yearExpense float64
	var yearIncome float64

	// Determine which months to show
	endMonth := 12
	if year == time.Now().Year() {
		endMonth = int(time.Now().Month())
	}

	// Format header
	msg.WriteString(fmt.Sprintf("📊 <b>%d Year Summary</b>\n\n", year))

	// Check if there are any transactions
	hasTransactions := false
	for m := 1; m <= endMonth; m++ {
		if _, ok := res[m]; ok {
			hasTransactions = true
			break
		}
	}

	if !hasTransactions {
		msg.WriteString(fmt.Sprintf("No transactions recorded for %d", year))
		return c.sendRecapWithNavigation(b, ctx, msg.String(), "year", year, 0)
	}

	// --- MONTHLY BREAKDOWN SECTION ---
	msg.WriteString("<b>Monthly Breakdown:</b>\n\n")

	for m := 1; m <= endMonth; m++ {
		monthT, hasTransactions := res[m]
		if !hasTransactions {
			continue // Skip months with no transactions
		}

		msg.WriteString(fmt.Sprintf("🗓 <b>%s</b>\n", time.Month(m).String()))
		var monthTotal float64

		if expenseAmount, ok := monthT[model.TypeExpense]; ok && expenseAmount > 0 {
			msg.WriteString(fmt.Sprintf("  💸 <b>Expenses:</b> %.2f€\n", expenseAmount))
			monthTotal -= expenseAmount
			yearExpense += expenseAmount
		}

		if incomeAmount, ok := monthT[model.TypeIncome]; ok && incomeAmount > 0 {
			msg.WriteString(fmt.Sprintf("  💰 <b>Income:</b> %.2f€\n", incomeAmount))
			monthTotal += incomeAmount
			yearIncome += incomeAmount
		}

		yearTotal += monthTotal

		var balanceEmoji string
		if monthTotal >= 0 {
			balanceEmoji = "✅"
		} else {
			balanceEmoji = "❌"
		}

		msg.WriteString(fmt.Sprintf("  %s <b>Balance:</b> %.2f€\n\n", balanceEmoji, monthTotal))
	}

	// --- YEAR TOTAL SECTION ---
	msg.WriteString("\n<b>💰 Year Summary</b>\n")

	// Add expense summary with category breakdown
	if yearExpense > 0 {
		msg.WriteString(fmt.Sprintf("💸 <b>Total Expenses:</b> %.2f€\n", yearExpense))

		// Add category breakdown for expenses
		if expenseCats, ok := categoryTotals[model.TypeExpense]; ok && len(expenseCats) > 0 {
			msg.WriteString("\n<b>Expense Categories:</b>\n")

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

			// Display top categories (limit to top 5 for readability)
			maxCategories := min(len(categories), 5)

			for i := range maxCategories {
				entry := categories[i]
				emoji := utils.GetCategoryEmoji(entry.Category)
				percentage := (entry.Amount / yearExpense) * 100
				msg.WriteString(fmt.Sprintf("  %s <b>%s:</b> %.2f€ (%.1f%%)\n",
					emoji, entry.Category, entry.Amount, percentage))
			}

			// Show "Other" for remaining categories if more than 5
			if len(categories) > maxCategories {
				var otherAmount float64
				for i := maxCategories; i < len(categories); i++ {
					otherAmount += categories[i].Amount
				}
				percentage := (otherAmount / yearExpense) * 100
				msg.WriteString(fmt.Sprintf("  📌 <b>Others:</b> %.2f€ (%.1f%%)\n",
					otherAmount, percentage))
			}

			msg.WriteString("\n")
		}
	}

	// Add income summary with category breakdown
	if yearIncome > 0 {
		msg.WriteString(fmt.Sprintf("💰 <b>Total Income:</b> %.2f€\n", yearIncome))

		// Add category breakdown for income
		if incomeCats, ok := categoryTotals[model.TypeIncome]; ok && len(incomeCats) > 0 {
			msg.WriteString("\n<b>Income Categories:</b>\n")

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

			// Display all income categories (usually fewer than expenses)
			for _, entry := range categories {
				emoji := utils.GetCategoryEmoji(entry.Category)
				percentage := (entry.Amount / yearIncome) * 100
				msg.WriteString(fmt.Sprintf("  %s <b>%s:</b> %.2f€ (%.1f%%)\n",
					emoji, entry.Category, entry.Amount, percentage))
			}

			msg.WriteString("\n")
		}
	}

	// Add final balance
	var balanceEmoji string
	if yearTotal >= 0 {
		balanceEmoji = "✅"
	} else {
		balanceEmoji = "❌"
	}

	msg.WriteString(fmt.Sprintf("\n%s <b>Year Balance:</b> %.2f€", balanceEmoji, yearTotal))

	return c.sendRecapWithNavigation(b, ctx, msg.String(), "year", year, 0)
	// return c.SendHomeKeyboard(b, ctx, msg.String())
}
