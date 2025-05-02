package client

import (
	"cashout/internal/model"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) YearRecap(b *gotgbot.Bot, ctx *ext.Context) error {
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

	// Get current year
	currentYear := time.Now().Year()

	// Get monthly totals for all months
	res, err := c.Repositories.Transactions.GetMonthlyTotalsCurrentYear(user.TgID)
	if err != nil {
		return err
	}

	// Get category breakdown for the entire year
	categoryTotals, err := c.Repositories.Transactions.GetYearCategorizedTotals(user.TgID, currentYear)
	if err != nil {
		return err
	}

	var msg strings.Builder
	var yearTotal float64
	var yearExpense float64
	var yearIncome float64

	currMonth := time.Now().Month()

	// Format header
	msg.WriteString(fmt.Sprintf("📊 <b>%d Year to Date Summary</b>\n\n", currentYear))

	// --- MONTHLY BREAKDOWN SECTION ---
	msg.WriteString("<b>Monthly Breakdown:</b>\n\n")

	for m := 1; m <= int(currMonth); m++ {
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
	msg.WriteString("\n<b>💰 Year to Date Summary</b>\n")

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
			maxCategories := 5
			if len(categories) < maxCategories {
				maxCategories = len(categories)
			}

			for i := 0; i < maxCategories; i++ {
				entry := categories[i]
				emoji := getCategoryEmoji(entry.Category)
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
				emoji := getCategoryEmoji(entry.Category)
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

	// Send the message
	ctx.EffectiveMessage.Reply(b, msg.String(), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
