package scheduler

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
)

func getNextFirstDayOfMonth() time.Time {
	now := time.Now()
	// Get first day of next month
	firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	return firstOfNextMonth
}

// createMonthlyReminders creates reminder records for all active users
func (s *Scheduler) createMonthlyReminders() error {
	s.logger.Info("Creating monthly reminders...")

	// Get all active users
	users, err := s.repositories.Reminders.GetAllActiveUsers()
	if err != nil {
		return fmt.Errorf("failed to get active users: %w", err)
	}

	// Schedule for 6:00 GMT on the first day of next month
	firstDay := getNextFirstDayOfMonth()
	scheduledFor := time.Date(firstDay.Year(), firstDay.Month(), firstDay.Day(), 6, 0, 0, 0, time.UTC)

	createdCount := 0
	for _, user := range users {
		err := s.repositories.Reminders.CreateOrUpdateMonthlyReminder(user.TgID, scheduledFor)
		if err != nil {
			s.logger.Errorf("Failed to create monthly reminder for user %d: %v", user.TgID, err)
			continue
		}
		createdCount++
	}

	s.logger.Infof("Created %d monthly reminders for %d active users", createdCount, len(users))
	return nil
}

// processMonthlyReminders processes all pending monthly reminders
func (s *Scheduler) processMonthlyReminders() error {
	// Get all pending reminders that should be sent now
	reminders, err := s.repositories.Reminders.GetPendingReminders(
		model.ReminderTypeMonthlyRecap,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to get pending monthly reminders: %w", err)
	}

	if len(reminders) == 0 {
		return nil
	}

	s.logger.Infof("Processing %d pending monthly reminders", len(reminders))

	for _, reminder := range reminders {
		// Update status to processing (with transaction to prevent double processing)
		err := s.repositories.Reminders.UpdateReminderStatusTransaction(
			reminder.ID,
			model.ReminderStatusProcessing,
			nil,
		)
		if err != nil {
			s.logger.Errorf("Failed to update monthly reminder %d to processing: %v", reminder.ID, err)
			continue
		}

		// Send the monthly recap
		err = s.sendMonthlyRecap(reminder.TgID)

		if err != nil {
			errMsg := err.Error()
			errors.Join(err, s.repositories.Reminders.UpdateReminderStatusTransaction(
				reminder.ID,
				model.ReminderStatusFailed,
				&errMsg,
			))
			s.logger.Errorf("Failed to send monthly recap for user %d: %v", reminder.TgID, err)
		} else {
			errors.Join(err, s.repositories.Reminders.UpdateReminderStatusTransaction(
				reminder.ID,
				model.ReminderStatusSent,
				nil,
			))
			s.logger.Infof("Successfully sent monthly recap to user %d", reminder.TgID)
		}
	}

	return nil
}

// sendMonthlyRecap sends the previous month's recap to a user
func (s *Scheduler) sendMonthlyRecap(tgID int64) error {
	// Get user data
	user, err := s.repositories.Users.GetByTgID(tgID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Calculate previous month
	now := time.Now().UTC()
	// Go to first day of current month
	firstOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	// Go back one day to get last day of previous month
	lastOfPrevMonth := firstOfCurrentMonth.AddDate(0, 0, -1)
	// Get year and month of previous month
	prevYear := lastOfPrevMonth.Year()
	prevMonth := int(lastOfPrevMonth.Month())

	// Get monthly totals
	totals, err := s.repositories.Transactions.GetMonthlyTotalsInYear(user.TgID, prevYear)
	if err != nil {
		return fmt.Errorf("failed to get monthly totals: %w", err)
	}

	// Get category breakdown
	categoryTotals, err := s.repositories.Transactions.GetMonthCategorizedTotals(user.TgID, prevYear, prevMonth)
	if err != nil {
		return fmt.Errorf("failed to get category totals: %w", err)
	}

	// Generate the recap message
	message := s.generateMonthlyRecapMessage(user, totals, categoryTotals, prevYear, prevMonth)

	// Send the message
	_, err = s.bot.SendMessage(tgID, message, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return err
}

// generateMonthlyRecapMessage generates the monthly recap message
func (s *Scheduler) generateMonthlyRecapMessage(user model.User, totals map[int]map[model.TransactionType]float64, categoryTotals map[model.TransactionType]map[model.TransactionCategory]float64, year int, month int) string {
	var text strings.Builder
	var monthTotal float64

	// Header
	text.WriteString(fmt.Sprintf("📅 <b>%s, here's your monthly recap!</b>\n\n", user.Name))
	text.WriteString(fmt.Sprintf("📊 <b>%s %d Summary</b>\n\n", time.Month(month).String(), year))

	t, ok := totals[month]
	if !ok {
		text.WriteString(fmt.Sprintf("You had no transactions in %s %d.\n\n", time.Month(month).String(), year))
		text.WriteString("💡 <i>Start tracking your expenses to get insights!</i>")
		return text.String()
	}

	// --- EXPENSES SECTION ---
	if expenseAmount, ok := t[model.TypeExpense]; ok && expenseAmount > 0 {
		monthTotal -= expenseAmount
		text.WriteString(fmt.Sprintf("💸 <b>Expenses:</b> %.2f€\n", expenseAmount))

		// Add top expense categories
		if expenseCats, ok := categoryTotals[model.TypeExpense]; ok && len(expenseCats) > 0 {
			text.WriteString("\n<b>Top Expenses:</b>\n")

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

			// Display top 5 categories
			limit := 5
			if len(categories) < limit {
				limit = len(categories)
			}

			for i := 0; i < limit; i++ {
				entry := categories[i]
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

		// Add income categories if multiple
		if incomeCats, ok := categoryTotals[model.TypeIncome]; ok && len(incomeCats) > 1 {
			text.WriteString("\n<b>Income Sources:</b>\n")

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

	text.WriteString(fmt.Sprintf("\n%s <b>Month Balance:</b> %.2f€\n", balanceEmoji, monthTotal))

	// --- AVERAGE DAILY SPENDING ---
	if expenseAmount, ok := t[model.TypeExpense]; ok && expenseAmount > 0 {
		daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
		avgDaily := expenseAmount / float64(daysInMonth)
		text.WriteString(fmt.Sprintf("📈 <b>Avg Daily Spending:</b> %.2f€\n", avgDaily))
	}

	// --- COMPARISON WITH PREVIOUS MONTH ---
	if month > 1 || year > 2020 {
		prevMonthNum := month - 1
		prevYearNum := year
		if prevMonthNum < 1 {
			prevMonthNum = 12
			prevYearNum--
		}

		if prevTotals, hasPrev := totals[prevMonthNum]; hasPrev && prevYearNum == year {
			text.WriteString("\n<b>vs Previous Month:</b>\n")

			// Compare expenses
			if currentExpense, ok := t[model.TypeExpense]; ok {
				if prevExpense, ok := prevTotals[model.TypeExpense]; ok && prevExpense > 0 {
					diff := currentExpense - prevExpense
					percentChange := (diff / prevExpense) * 100

					if diff > 0 {
						text.WriteString(fmt.Sprintf("  📈 Expenses: +%.2f€ (+%.1f%%)\n", diff, percentChange))
					} else {
						text.WriteString(fmt.Sprintf("  📉 Expenses: %.2f€ (%.1f%%)\n", diff, percentChange))
					}
				}
			}

			// Compare income
			if currentIncome, ok := t[model.TypeIncome]; ok {
				if prevIncome, ok := prevTotals[model.TypeIncome]; ok && prevIncome > 0 {
					diff := currentIncome - prevIncome
					percentChange := (diff / prevIncome) * 100

					if diff > 0 {
						text.WriteString(fmt.Sprintf("  📈 Income: +%.2f€ (+%.1f%%)\n", diff, percentChange))
					} else {
						text.WriteString(fmt.Sprintf("  📉 Income: %.2f€ (%.1f%%)\n", diff, percentChange))
					}
				}
			}
		}
	}

	text.WriteString("\n💡 <i>Type /month to see this month's progress!</i>")

	return text.String()
}
