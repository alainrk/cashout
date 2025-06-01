package db

import (
	"cashout/internal/model"
	"time"

	"gorm.io/gorm"
)

// CreateReminder creates a new reminder record
func (db *DB) CreateReminder(reminder *model.Reminder) error {
	return db.conn.Create(reminder).Error
}

// GetPendingReminders retrieves all pending reminders that should be processed
func (db *DB) GetPendingReminders(reminderType model.ReminderType, scheduledBefore time.Time) ([]model.Reminder, error) {
	var reminders []model.Reminder
	result := db.conn.Where("type = ? AND status = ? AND scheduled_for <= ?",
		reminderType, model.ReminderStatusPending, scheduledBefore).
		Find(&reminders)

	if result.Error != nil {
		return nil, result.Error
	}
	return reminders, nil
}

// UpdateReminderStatusTransaction updates a reminder's status within a transaction
// This ensures consistency and prevents double processing
func (db *DB) UpdateReminderStatusTransaction(reminderID int64, status model.ReminderStatus, errorMsg *string) error {
	return db.conn.Transaction(func(tx *gorm.DB) error {
		// First, check if the reminder is still pending
		var reminder model.Reminder
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ? AND status = ?", reminderID, model.ReminderStatusPending).
			First(&reminder)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Reminder already processed by another instance
				return nil
			}
			return result.Error
		}

		// Update the reminder
		updates := map[string]interface{}{
			"status": status,
		}

		if status == model.ReminderStatusSent || status == model.ReminderStatusFailed {
			now := time.Now()
			updates["processed_at"] = &now
		}

		if errorMsg != nil {
			updates["error_message"] = errorMsg
		}

		return tx.Model(&model.Reminder{}).
			Where("id = ?", reminderID).
			Updates(updates).Error
	})
}

// CreateOrUpdateWeeklyReminder creates or updates a weekly reminder for a user
func (db *DB) CreateOrUpdateWeeklyReminder(tgID int64, scheduledFor time.Time) error {
	reminder := model.Reminder{
		TgID:         tgID,
		Type:         model.ReminderTypeWeeklyRecap,
		Status:       model.ReminderStatusPending,
		ScheduledFor: scheduledFor,
	}

	// Use ON CONFLICT to update if exists
	result := db.conn.Exec(`
		INSERT INTO reminders (tg_id, type, status, scheduled_for, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (tg_id, type, scheduled_for) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			updated_at = CURRENT_TIMESTAMP
		WHERE reminders.status != ?
	`, reminder.TgID, reminder.Type, reminder.Status, reminder.ScheduledFor, model.ReminderStatusSent)

	return result.Error
}

// CreateOrUpdateMonthlyReminder creates or updates a monthly reminder for a user
func (db *DB) CreateOrUpdateMonthlyReminder(tgID int64, scheduledFor time.Time) error {
	reminder := model.Reminder{
		TgID:         tgID,
		Type:         model.ReminderTypeMonthlyRecap,
		Status:       model.ReminderStatusPending,
		ScheduledFor: scheduledFor,
	}

	// Use ON CONFLICT to update if exists
	result := db.conn.Exec(`
		INSERT INTO reminders (tg_id, type, status, scheduled_for, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (tg_id, type, scheduled_for) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			updated_at = CURRENT_TIMESTAMP
		WHERE reminders.status != ?
	`, reminder.TgID, reminder.Type, reminder.Status, reminder.ScheduledFor, model.ReminderStatusSent)

	return result.Error
}

// GetAllActiveUsers
func (db *DB) GetAllActiveUsers() ([]model.User, error) {
	var users []model.User
	result := db.conn.Distinct("users.*").
		// TODO: Find something else here, for now it's not a problem
		// Joins("JOIN transactions ON users.tg_id = transactions.tg_id").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
