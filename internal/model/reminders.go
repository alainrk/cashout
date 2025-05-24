package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

// ReminderType represents the type of reminder
type ReminderType string

// Reminder types
const (
	ReminderTypeWeeklyRecap  ReminderType = "weekly_recap"
	ReminderTypeMonthlyRecap ReminderType = "monthly_recap"
	ReminderTypeYearlyRecap  ReminderType = "yearly_recap"
)

// Value implements the driver.Valuer interface for ReminderType
func (r ReminderType) Value() (driver.Value, error) {
	return string(r), nil
}

// Scan implements the sql.Scanner interface for ReminderType
func (r *ReminderType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("reminder type cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid reminder type")
	}

	*r = ReminderType(strVal)
	return nil
}

// ReminderStatus represents the status of a reminder
type ReminderStatus string

// Reminder statuses
const (
	ReminderStatusPending    ReminderStatus = "pending"
	ReminderStatusProcessing ReminderStatus = "processing"
	ReminderStatusSent       ReminderStatus = "sent"
	ReminderStatusFailed     ReminderStatus = "failed"
)

// Value implements the driver.Valuer interface for ReminderStatus
func (s ReminderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the sql.Scanner interface for ReminderStatus
func (s *ReminderStatus) Scan(value interface{}) error {
	if value == nil {
		return errors.New("reminder status cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid reminder status")
	}

	*s = ReminderStatus(strVal)
	return nil
}

// Reminder represents the reminders table structure
type Reminder struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement"`
	TgID         int64          `gorm:"column:tg_id;not null;index"`
	Type         ReminderType   `gorm:"column:type;not null;type:reminder_type;index"`
	Status       ReminderStatus `gorm:"column:status;not null;type:reminder_status;default:'pending';index"`
	ScheduledFor time.Time      `gorm:"column:scheduled_for;not null;index"`
	ProcessedAt  *time.Time     `gorm:"column:processed_at"`
	ErrorMessage *string        `gorm:"column:error_message;type:text"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`

	// Association to User (optional)
	User *User `gorm:"foreignKey:TgID;references:TgID"`
}

// TableName overrides the table name
func (Reminder) TableName() string {
	return "reminders"
}
