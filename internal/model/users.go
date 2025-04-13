package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StateType represents the status of the user in the session
type StateType string

// State type constants
const (
	StateStart   StateType = "start"   // When the user has just started the conversation with the bot
	StateNormal  StateType = "normal"  // When the user can send commands or a transaction "out-of-the-blue" (ootb)
	StateWaiting StateType = "waiting" // We are waiting a followup response from the user
)

// CommandType represents the type of command sent by the user
type CommandType string

// Command type constants
const (
	CommandNone             CommandType = "none"
	CommandStart            CommandType = "start"
	CommandHelp             CommandType = "help"
	CommandCancel           CommandType = "cancel"
	CommandConfirm          CommandType = "confirm"
	CommandAddIncomeIntent  CommandType = "add_income_intent"
	CommandAddExpenseIntent CommandType = "add_expense_intent"
	CommandAddTransaction   CommandType = "add_transaction" // User has sent their transaction text
)

// User represents the users table structure
type User struct {
	TgID        int64       `gorm:"column:tg_id;primaryKey"`
	TgUsername  string      `gorm:"column:tg_username;unique"`
	TgFirstname string      `gorm:"column:tg_firstname"`
	TgLastname  string      `gorm:"column:tg_lastname"`
	Name        string      `gorm:"column:name;name"`
	Session     UserSession `gorm:"column:session;type:jsonb"`
	Settings    JSONData    `gorm:"column:settings;type:jsonb"`
	CreatedAt   time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;autoUpdateTime"`
}

type UserSession struct {
	Iterations  uint64      `json:"iterations"`
	State       StateType   `json:"state"`
	LastCommand CommandType `json:"last_command"`
	LastMessage string      `json:"last_message"`
	Body        string      `json:"body"`
}

// Value makes the UserSession struct implement the driver.Valuer interface
func (s UserSession) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan makes the UserSession struct implement the sql.Scanner interface
func (s *UserSession) Scan(value interface{}) error {
	if value == nil {
		*s = UserSession{State: "", LastCommand: CommandNone, LastMessage: ""}
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
