package db

import (
	"time"
)

// User represents the users table structure
type User struct {
	TgID        int64     `gorm:"column:tg_id;primaryKey"`
	TgUsername  string    `gorm:"column:tg_username;unique"`
	TgFirstname string    `gorm:"column:tg_firstname"`
	TgLastname  string    `gorm:"column:tg_lastname"`
	Name        string    `gorm:"column:name;name"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}
