package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

// TransactionCategory represents the category of an transaction or income
type TransactionCategory string

// Transaction categories
const (
	CategorySalary        TransactionCategory = "Salary"
	CategoryOtherIncomes  TransactionCategory = "OtherIncomes"
	CategoryCar           TransactionCategory = "Car"
	CategoryClothes       TransactionCategory = "Clothes"
	CategoryGrocery       TransactionCategory = "Grocery"
	CategoryHouse         TransactionCategory = "House"
	CategoryBills         TransactionCategory = "Bills"
	CategoryEntertainment TransactionCategory = "Entertainment"
	CategorySport         TransactionCategory = "Sport"
	CategoryEatingOut     TransactionCategory = "EatingOut"
	CategoryTransport     TransactionCategory = "Transport"
	CategoryLearning      TransactionCategory = "Learning"
	CategoryToiletry      TransactionCategory = "Toiletry"
	CategoryHealth        TransactionCategory = "Health"
	CategoryTech          TransactionCategory = "Tech"
	CategoryGifts         TransactionCategory = "Gifts"
	CategoryTravel        TransactionCategory = "Travel"
	CategoryOtherExpenses TransactionCategory = "OtherExpenses"
)

func IsValidTransactionCategory(category string) bool {
	for _, c := range GetTransactionCategories() {
		if c == category {
			return true
		}
	}
	return false
}

// Value implements the driver.Valuer interface for TransactionCategory
func (t TransactionCategory) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the sql.Scanner interface for TransactionCategory
func (t *TransactionCategory) Scan(value interface{}) error {
	if value == nil {
		return errors.New("transaction category cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid transaction category type")
	}

	*t = TransactionCategory(strVal)
	return nil
}

// TransactionType represents whether the entry is an income or an transaction
type TransactionType string

// Transaction types
const (
	TypeIncome  TransactionType = "Income"
	TypeExpense TransactionType = "Expense"
)

// Value implements the driver.Valuer interface for TransactionType
func (t TransactionType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the sql.Scanner interface for TransactionType
func (t *TransactionType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("transaction type cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid transaction type")
	}

	*t = TransactionType(strVal)
	return nil
}

// CurrencyType represents the currency of an transaction or income
type CurrencyType string

// Currency types
const (
	CurrencyEUR CurrencyType = "EUR"
	CurrencyUSD CurrencyType = "USD"
	CurrencyGBP CurrencyType = "GBP"
	CurrencyJPY CurrencyType = "JPY"
	CurrencyCHF CurrencyType = "CHF"
)

// Value implements the driver.Valuer interface for CurrencyType
func (t CurrencyType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the sql.Scanner interface for CurrencyType
func (t *CurrencyType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("currency type cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid currency type")
	}

	*t = CurrencyType(strVal)
	return nil
}

// Transaction represents the transactions table structure
type Transaction struct {
	ID          int64               `gorm:"column:id;primaryKey;autoIncrement"`
	TgID        int64               `gorm:"column:tg_id;not null;index"`
	Date        time.Time           `gorm:"column:date;not null;type:date;default:CURRENT_DATE;index"`
	Type        TransactionType     `gorm:"column:type;not null;type:transaction_type;index"`
	Category    TransactionCategory `gorm:"column:category;not null;type:transaction_category;index"`
	Amount      float64             `gorm:"column:amount;not null;type:decimal(15,2)"`
	Currency    CurrencyType        `gorm:"column:currency;not null;type:currency_type;default:'EUR'"`
	Description string              `gorm:"column:description;type:text"`
	CreatedAt   time.Time           `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time           `gorm:"column:updated_at;autoUpdateTime"`

	// Association to User (optional)
	User *User `gorm:"foreignKey:TgID;references:TgID"`
}

// TableName overrides the table name
func (Transaction) TableName() string {
	return "transactions"
}

// Get category enum values as a slice of strings
func GetTransactionCategories() []string {
	return []string{
		string(CategorySalary),
		string(CategoryOtherIncomes),
		string(CategoryCar),
		string(CategoryClothes),
		string(CategoryGrocery),
		string(CategoryHouse),
		string(CategoryBills),
		string(CategoryEntertainment),
		string(CategorySport),
		string(CategoryEatingOut),
		string(CategoryTransport),
		string(CategoryLearning),
		string(CategoryToiletry),
		string(CategoryHealth),
		string(CategoryTech),
		string(CategoryGifts),
		string(CategoryTravel),
		string(CategoryOtherExpenses),
	}
}

// Get transaction type enum values as a slice of strings
func GetTransactionTypes() []string {
	return []string{
		string(TypeIncome),
		string(TypeExpense),
	}
}

// Get currency type enum values as a slice of strings
func GetCurrencyTypes() []string {
	return []string{
		string(CurrencyEUR),
		string(CurrencyUSD),
		string(CurrencyGBP),
		string(CurrencyJPY),
		string(CurrencyCHF),
	}
}
