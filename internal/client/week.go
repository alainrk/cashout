// internal/client/week.go
package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// WeekRecap returns to the user the breakdown and the total for the expenses and income of the current week.
func (c *Client) WeekRecap(b *gotgbot.Bot, ctx *ext.Context) error {
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

	// Get current week boundaries (Monday to Sunday)
	now := time.Now()
	weekday := int(now.Weekday())
	// If Sunday (0), make it 7 for calculation
	if weekday == 0 {
		weekday = 7
	}
	// Calculate days to subtract to get to Monday
	daysToMonday := weekday - 1
	startOfWeek := now.AddDate(0, 0, -daysToMonday)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	// End of week is Sunday (6 days after Monday)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 999999999, endOfWeek.Location())

	// Get transactions for the week
	transactions, err := c.Repositories.Transactions.GetUserTransactionsByDateRange(user.TgID, startOfWeek, endOfWeek)
	if err != nil {
		return fmt.Errorf("failed to get weekly transactions: %w", err)
	}

	if len(transactions) == 0 {
		txt := fmt.Sprintf("No transactions for this week (%s - %s)",
			startOfWeek.Format("02 Jan"),
			endOfWeek.Format("02 Jan"))
		return c.SendHomeKeyboard(b, ctx, txt)
	}

	// Calculate totals by type and category
	typeTotals := make(map[model.TransactionType]float64)
	categoryTotals := make(map[model.TransactionType]map[model.TransactionCategory]float64)
	dailyTotals := make(map[string]map[model.TransactionType]float64)

	// Initialize category totals map
	categoryTotals[model.TypeExpense] = make(map[model.TransactionCategory]float64)
	categoryTotals[model.TypeIncome] = make(map[model.TransactionCategory]float64)

	for _, t := range transactions {
		// Type totals
		typeTotals[t.Type] += t.Amount

		// Category totals
		categoryTotals[t.Type][t.Category] += t.Amount

		// Daily totals
		dayKey := t.Date.Format("Mon 02")
		if dailyTotals[dayKey] == nil {
			dailyTotals[dayKey] = make(map[model.TransactionType]float64)
		}
		dailyTotals[dayKey][t.Type] += t.Amount
	}

	// Format the message
	var text strings.Builder
	var weekTotal float64

	// Header with week dates
	text.WriteString(fmt.Sprintf("📊 <b>Week %s - %s</b>\n\n",
		startOfWeek.Format("02 Jan"),
		endOfWeek.Format("02 Jan")))

	// --- DAILY BREAKDOWN ---
	text.WriteString("<b>Daily Activity:</b>\n")

	// Get all days in order
	var dayKeys []string
	for day := startOfWeek; !day.After(endOfWeek); day = day.AddDate(0, 0, 1) {
		dayKey := day.Format("Mon 02")
		dayKeys = append(dayKeys, dayKey)
	}

	hasActivity := false
	for _, dayKey := range dayKeys {
		if totals, exists := dailyTotals[dayKey]; exists {
			hasActivity = true
			dayBalance := 0.0

			text.WriteString(fmt.Sprintf("\n📅 <b>%s</b>\n", dayKey))

			if expense, ok := totals[model.TypeExpense]; ok && expense > 0 {
				text.WriteString(fmt.Sprintf("  💸 %.2f€\n", expense))
				dayBalance -= expense
			}

			if income, ok := totals[model.TypeIncome]; ok && income > 0 {
				text.WriteString(fmt.Sprintf("  💰 %.2f€\n", income))
				dayBalance += income
			}

			// Show daily balance only if there were both income and expenses
			if len(totals) > 1 {
				emoji := "✅"
				if dayBalance < 0 {
					emoji = "❌"
				}
				text.WriteString(fmt.Sprintf("  %s Balance: %.2f€\n", emoji, dayBalance))
			}
		}
	}

	if !hasActivity {
		text.WriteString("\nNo transactions recorded this week.\n")
	}

	text.WriteString("\n")

	// --- EXPENSES SECTION ---
	if expenseAmount, ok := typeTotals[model.TypeExpense]; ok && expenseAmount > 0 {
		weekTotal -= expenseAmount
		text.WriteString(fmt.Sprintf("💸 <b>Total Expenses:</b> %.2f€\n", expenseAmount))

		// Add category breakdown for expenses
		if expenseCats := categoryTotals[model.TypeExpense]; len(expenseCats) > 0 {
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
	if incomeAmount, ok := typeTotals[model.TypeIncome]; ok && incomeAmount > 0 {
		weekTotal += incomeAmount
		text.WriteString(fmt.Sprintf("💰 <b>Total Income:</b> %.2f€\n", incomeAmount))

		// Add category breakdown for income
		if incomeCats := categoryTotals[model.TypeIncome]; len(incomeCats) > 0 {
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
	if weekTotal >= 0 {
		balanceEmoji = "✅"
	} else {
		balanceEmoji = "❌"
	}

	text.WriteString(fmt.Sprintf("\n%s <b>Week Balance:</b> %.2f€", balanceEmoji, weekTotal))

	// --- AVERAGE DAILY SPENDING ---
	if expenseAmount, ok := typeTotals[model.TypeExpense]; ok && expenseAmount > 0 {
		avgDaily := expenseAmount / 7
		text.WriteString(fmt.Sprintf("\n📈 <b>Avg Daily Spending:</b> %.2f€", avgDaily))
	}

	return c.SendHomeKeyboard(b, ctx, text.String())
}
