package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StateType represents the last state of the user (if any)
type StateType string

// State type constants
const (
	// Normal state where the user has to give commands.
	StateNormal StateType = "normal"
	// Only when the user has started the bot.
	StateStart StateType = "start"
	// The user has to add an expense after having set the bot to accept it.
	StateInsertingExpense StateType = "inserting_expense"
	// The user has to add an income after having set the bot to accept it.
	StateInsertingIncome StateType = "inserting_income"
	// The user has to edit the transaction, during an insert flow
	StateEditingTransactionDate     StateType = "editing_transaction_date"
	StateEditingTransactionCategory StateType = "editing_transaction_category"
	StateEditingTransactionAmount   StateType = "editing_transaction_amount"
	// The user has to edit the transaction, during an edit flow
	StateTopLevelEditingTransactionDate     StateType = "top_level_editing_transaction_date"
	StateTopLevelEditingTransactionCategory StateType = "top_level_editing_transaction_category"
	StateTopLevelEditingTransactionAmount   StateType = "top_level_editing_transaction_amount"
	// The user has to confirm an action.
	StateWaitingConfirm StateType = "waiting_confirm"
)

// CommandType represents the type of command sent by the user
type CommandType string

// User represents the users table structure
type User struct {
	TgID        int64       `gorm:"column:tg_id;primaryKey"`
	TgUsername  string      `gorm:"column:tg_username;unique"`
	TgFirstname string      `gorm:"column:tg_firstname"`
	TgLastname  string      `gorm:"column:tg_lastname"`
	Name        string      `gorm:"column:name;name"`
	Session     UserSession `gorm:"column:session;type:jsonb"`
	CreatedAt   time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;autoUpdateTime"`
}

type UserSession struct {
	State       StateType `json:"state"`
	LastMessage string    `json:"last_message"`
	Body        string    `json:"body"`
}

// Value makes the UserSession struct implement the driver.Valuer interface
func (s UserSession) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan makes the UserSession struct implement the sql.Scanner interface
func (s *UserSession) Scan(value interface{}) error {
	if value == nil {
		*s = UserSession{State: "", LastMessage: ""}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, s)
}

// JSONData is a custom type for handling JSON in the database
type JSONData map[string]interface{}

// Value makes the JSONData struct implement the driver.Valuer interface
func (j JSONData) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan makes the JSONData struct implement the sql.Scanner interface
func (j *JSONData) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONData)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}
