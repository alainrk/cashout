package client

import (
	"fmt"
	"happypoor/internal/model"
	"sort"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// MonthRecap returns to the user the breakdown and the total for the expenses and income of the current month.
func (c *Client) MonthRecap(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	// Get current month and year
	now := time.Now()
	month := now.Month()
	year := now.Year()

	// Get monthly totals
	totals, err := c.Repositories.Transactions.GetMonthlyTotalsCurrentYear(user.TgID)
	if err != nil {
		return err
	}

	// Get category breakdown
	categoryTotals, err := c.Repositories.Transactions.GetMonthCategorizedTotals(user.TgID, year, int(month))
	if err != nil {
		return err
	}

	t, ok := totals[int(month)]
	if !ok {
		ctx.EffectiveMessage.Reply(b, "No transactions for this month", &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Format the message
	var msg strings.Builder
	var monthTotal float64

	// Header with month name
	msg.WriteString(fmt.Sprintf("📊 <b>%s %d Summary</b>\n\n", month.String(), year))

	// --- EXPENSES SECTION ---
	if expenseAmount, ok := t[model.TypeExpense]; ok && expenseAmount > 0 {
		monthTotal -= expenseAmount
		msg.WriteString(fmt.Sprintf("💸 <b>Expenses:</b> %.2f€\n", expenseAmount))

		// Add category breakdown for expenses
		if expenseCats, ok := categoryTotals[model.TypeExpense]; ok && len(expenseCats) > 0 {
			msg.WriteString("\n<b>Expense Breakdown:</b>\n")

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
				emoji := getCategoryEmoji(entry.Category)
				percentage := (entry.Amount / expenseAmount) * 100
				msg.WriteString(fmt.Sprintf("  %s <b>%s:</b> %.2f€ (%.1f%%)\n",
					emoji, entry.Category, entry.Amount, percentage))
			}
			msg.WriteString("\n")
		}
	}

	// --- INCOME SECTION ---
	if incomeAmount, ok := t[model.TypeIncome]; ok && incomeAmount > 0 {
		monthTotal += incomeAmount
		msg.WriteString(fmt.Sprintf("💰 <b>Income:</b> %.2f€\n", incomeAmount))

		// Add category breakdown for income
		if incomeCats, ok := categoryTotals[model.TypeIncome]; ok && len(incomeCats) > 0 {
			msg.WriteString("\n<b>Income Breakdown:</b>\n")

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
				emoji := getCategoryEmoji(entry.Category)
				percentage := (entry.Amount / incomeAmount) * 100
				msg.WriteString(fmt.Sprintf("  %s <b>%s:</b> %.2f€ (%.1f%%)\n",
					emoji, entry.Category, entry.Amount, percentage))
			}
			msg.WriteString("\n")
		}
	}

	// --- TOTAL BALANCE ---
	var balanceEmoji string
	if monthTotal >= 0 {
		balanceEmoji = "✅"
	} else {
		balanceEmoji = "❌"
	}

	msg.WriteString(fmt.Sprintf("\n%s <b>Month Balance:</b> %.2f€", balanceEmoji, monthTotal))

	// Send the message
	ctx.EffectiveMessage.Reply(b, msg.String(), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
