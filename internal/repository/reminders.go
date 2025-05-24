package repository

import (
	"cashout/internal/model"
	"time"
)

type Reminders struct {
	Repository
}

func (r *Reminders) CreateReminder(reminder *model.Reminder) error {
	return r.DB.CreateReminder(reminder)
}

func (r *Reminders) GetPendingReminders(reminderType model.ReminderType, scheduledBefore time.Time) ([]model.Reminder, error) {
	return r.DB.GetPendingReminders(reminderType, scheduledBefore)
}

func (r *Reminders) UpdateReminderStatusTransaction(reminderID int64, status model.ReminderStatus, errorMsg *string) error {
	return r.DB.UpdateReminderStatusTransaction(reminderID, status, errorMsg)
}

func (r *Reminders) CreateOrUpdateWeeklyReminder(tgID int64, scheduledFor time.Time) error {
	return r.DB.CreateOrUpdateWeeklyReminder(tgID, scheduledFor)
}

func (r *Reminders) GetAllActiveUsers() ([]model.User, error) {
	return r.DB.GetAllActiveUsers()
}
