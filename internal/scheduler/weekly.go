package scheduler

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
)

func getNextMonday() time.Time {
	now := time.Now()
	daysUntilMonday := (8 - int(now.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	nextMonday := now.AddDate(0, 0, daysUntilMonday)
	return nextMonday
}

// createWeeklyReminders creates reminder records for all active users
func (s *Scheduler) createWeeklyReminders() error {
	s.logger.Info("Creating weekly reminders...")

	// Get all active users
	users, err := s.repositories.Reminders.GetAllActiveUsers()
	if err != nil {
		return fmt.Errorf("failed to get active users: %w", err)
	}

	// Schedule for 6:00 GMT on Monday
	monday := getNextMonday()
	scheduledFor := time.Date(monday.Year(), monday.Month(), monday.Day(), 6, 0, 0, 0, time.UTC)

	createdCount := 0
	for _, user := range users {
		err := s.repositories.Reminders.CreateOrUpdateWeeklyReminder(user.TgID, scheduledFor)
		if err != nil {
			s.logger.Errorf("Failed to create reminder for user %d: %v", user.TgID, err)
			continue
		}
		createdCount++
	}

	s.logger.Infof("Created %d weekly reminders for %d active users", createdCount, len(users))
	return nil
}

// processWeeklyReminders processes all pending weekly reminders
func (s *Scheduler) processWeeklyReminders() error {
	// Get all pending reminders that should be sent now
	reminders, err := s.repositories.Reminders.GetPendingReminders(
		model.ReminderTypeWeeklyRecap,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to get pending reminders: %w", err)
	}

	if len(reminders) == 0 {
		return nil
	}

	s.logger.Infof("Processing %d pending weekly reminders", len(reminders))

	for _, reminder := range reminders {
		// Update status to processing (with transaction to prevent double processing)
		err := s.repositories.Reminders.UpdateReminderStatusTransaction(
			reminder.ID,
			model.ReminderStatusProcessing,
			nil,
		)
		if err != nil {
			s.logger.Errorf("Failed to update reminder %d to processing: %v", reminder.ID, err)
			continue
		}

		// Send the weekly recap
		err = s.sendWeeklyRecap(reminder.TgID)

		if err != nil {
			errMsg := err.Error()
			errors.Join(err, s.repositories.Reminders.UpdateReminderStatusTransaction(
				reminder.ID,
				model.ReminderStatusFailed,
				&errMsg,
			))
			s.logger.Errorf("Failed to send weekly recap for user %d: %v", reminder.TgID, err)
		} else {
			errors.Join(err, s.repositories.Reminders.UpdateReminderStatusTransaction(
				reminder.ID,
				model.ReminderStatusSent,
				nil,
			))
			s.logger.Infof("Successfully sent weekly recap to user %d", reminder.TgID)
		}
	}

	return nil
}

// sendWeeklyRecap sends the previous week's recap to a user
func (s *Scheduler) sendWeeklyRecap(tgID int64) error {
	// Get user data
	user, err := s.repositories.Users.GetByTgID(tgID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Calculate previous week boundaries (Monday to Sunday)
	now := time.Now().UTC()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	// Get to this week's Monday
	daysToMonday := weekday - 1
	thisMonday := now.AddDate(0, 0, -daysToMonday)

	// Previous week is 7 days before
	startOfPrevWeek := thisMonday.AddDate(0, 0, -7)
	startOfPrevWeek = time.Date(startOfPrevWeek.Year(), startOfPrevWeek.Month(), startOfPrevWeek.Day(), 0, 0, 0, 0, time.UTC)

	endOfPrevWeek := startOfPrevWeek.AddDate(0, 0, 6)
	endOfPrevWeek = time.Date(endOfPrevWeek.Year(), endOfPrevWeek.Month(), endOfPrevWeek.Day(), 23, 59, 59, 999999999, time.UTC)

	// Get transactions for the previous week
	transactions, err := s.repositories.Transactions.GetUserTransactionsByDateRange(user.TgID, startOfPrevWeek, endOfPrevWeek)
	if err != nil {
		return fmt.Errorf("failed to get weekly transactions: %w", err)
	}

	// Generate the recap message
	message := s.generateWeeklyRecapMessage(user, transactions, startOfPrevWeek, endOfPrevWeek)

	// Send the message
	_, err = s.bot.SendMessage(tgID, message, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return err
}

// generateWeeklyRecapMessage generates the weekly recap message
// This reuses the logic from the WeekRecap function but adapted for previous week
func (s *Scheduler) generateWeeklyRecapMessage(user model.User, transactions []model.Transaction, startOfWeek, endOfWeek time.Time) string {
	var text strings.Builder

	// Header
	text.WriteString(fmt.Sprintf("🗓 <b>%s, here's your weekly recap!</b>\n\n", user.Name))
	text.WriteString(fmt.Sprintf("📊 <b>Week %s - %s</b>\n\n",
		startOfWeek.Format("02 Jan"),
		endOfWeek.Format("02 Jan")))

	if len(transactions) == 0 {
		text.WriteString("You had no transactions last week.\n\n")
		text.WriteString("💡 <i>Start tracking your expenses to get insights!</i>")
		return text.String()
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

	var weekTotal float64

	// Summary section
	if expenseAmount, ok := typeTotals[model.TypeExpense]; ok && expenseAmount > 0 {
		weekTotal -= expenseAmount
		text.WriteString(fmt.Sprintf("💸 <b>Total Expenses:</b> %.2f€\n", expenseAmount))
	}

	if incomeAmount, ok := typeTotals[model.TypeIncome]; ok && incomeAmount > 0 {
		weekTotal += incomeAmount
		text.WriteString(fmt.Sprintf("💰 <b>Total Income:</b> %.2f€\n", incomeAmount))
	}

	// Balance
	var balanceEmoji string
	if weekTotal >= 0 {
		balanceEmoji = "✅"
	} else {
		balanceEmoji = "❌"
	}
	text.WriteString(fmt.Sprintf("\n%s <b>Week Balance:</b> %.2f€\n", balanceEmoji, weekTotal))

	// Top expense categories (if any)
	if expenseCats := categoryTotals[model.TypeExpense]; len(expenseCats) > 0 {
		text.WriteString("\n<b>Top Expenses:</b>\n")

		// Get top 3 categories
		type catAmount struct {
			cat    model.TransactionCategory
			amount float64
		}
		var sorted []catAmount
		for cat, amount := range expenseCats {
			sorted = append(sorted, catAmount{cat, amount})
		}
		// Sort by amount descending
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[j].amount > sorted[i].amount {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}

		// Show top 3
		limit := 3
		if len(sorted) < limit {
			limit = len(sorted)
		}
		for i := 0; i < limit; i++ {
			emoji := utils.GetCategoryEmoji(sorted[i].cat)
			text.WriteString(fmt.Sprintf("  %s %s: %.2f€\n", emoji, sorted[i].cat, sorted[i].amount))
		}
	}

	// Average daily spending
	if expenseAmount, ok := typeTotals[model.TypeExpense]; ok && expenseAmount > 0 {
		avgDaily := expenseAmount / 7
		text.WriteString(fmt.Sprintf("\n📈 <b>Avg Daily Spending:</b> %.2f€\n", avgDaily))
	}

	text.WriteString("\n💡 <i>Type /week to see this week's progress!</i>")

	return text.String()
}
