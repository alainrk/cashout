package db

import (
	"database/sql/driver"
	"errors"
	"time"
)

// ExpenseCategory represents the category of an expense or income
type ExpenseCategory string

// Expense categories
const (
	CategorySalary        ExpenseCategory = "Salary"
	CategoryOtherIncomes  ExpenseCategory = "OtherIncomes"
	CategoryCar           ExpenseCategory = "Car"
	CategoryClothes       ExpenseCategory = "Clothes"
	CategoryGrocery       ExpenseCategory = "Grocery"
	CategoryHouse         ExpenseCategory = "House"
	CategoryBills         ExpenseCategory = "Bills"
	CategoryEntertainment ExpenseCategory = "Entertainment"
	CategorySport         ExpenseCategory = "Sport"
	CategoryEatingOut     ExpenseCategory = "EatingOut"
	CategoryTransport     ExpenseCategory = "Transport"
	CategoryLearning      ExpenseCategory = "Learning"
	CategoryToiletry      ExpenseCategory = "Toiletry"
	CategoryHealth        ExpenseCategory = "Health"
	CategoryTech          ExpenseCategory = "Tech"
	CategoryGifts         ExpenseCategory = "Gifts"
	CategoryTravel        ExpenseCategory = "Travel"
)

// Value implements the driver.Valuer interface for ExpenseCategory
func (ec ExpenseCategory) Value() (driver.Value, error) {
	return string(ec), nil
}

// Scan implements the sql.Scanner interface for ExpenseCategory
func (ec *ExpenseCategory) Scan(value interface{}) error {
	if value == nil {
		return errors.New("expense category cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid expense category type")
	}

	*ec = ExpenseCategory(strVal)
	return nil
}

// ExpenseType represents whether the entry is an income or an expense
type ExpenseType string

// Expense types
const (
	TypeIncome  ExpenseType = "Income"
	TypeExpense ExpenseType = "Expense"
)

// Value implements the driver.Valuer interface for ExpenseType
func (et ExpenseType) Value() (driver.Value, error) {
	return string(et), nil
}

// Scan implements the sql.Scanner interface for ExpenseType
func (et *ExpenseType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("expense type cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid expense type")
	}

	*et = ExpenseType(strVal)
	return nil
}

// CurrencyType represents the currency of an expense or income
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
func (ct CurrencyType) Value() (driver.Value, error) {
	return string(ct), nil
}

// Scan implements the sql.Scanner interface for CurrencyType
func (ct *CurrencyType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("currency type cannot be null")
	}

	strVal, ok := value.(string)
	if !ok {
		return errors.New("invalid currency type")
	}

	*ct = CurrencyType(strVal)
	return nil
}

// Expense represents the expenses table structure
type Expense struct {
	ID          int64           `gorm:"column:id;primaryKey;autoIncrement"`
	TgID        int64           `gorm:"column:tg_id;not null;index"`
	Date        time.Time       `gorm:"column:date;not null;type:date;default:CURRENT_DATE;index"`
	Type        ExpenseType     `gorm:"column:type;not null;type:expense_type;index"`
	Category    ExpenseCategory `gorm:"column:category;not null;type:expense_category;index"`
	Amount      float64         `gorm:"column:amount;not null;type:decimal(15,2)"`
	Currency    CurrencyType    `gorm:"column:currency;not null;type:currency_type;default:'EUR'"`
	Description string          `gorm:"column:description;type:text"`
	CreatedAt   time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;autoUpdateTime"`

	// Association to User (optional)
	User *User `gorm:"foreignKey:TgID;references:TgID"`
}

// TableName overrides the table name
func (Expense) TableName() string {
	return "expenses"
}

// Get category enum values as a slice of strings
func GetExpenseCategories() []string {
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
	}
}

// Get expense type enum values as a slice of strings
func GetExpenseTypes() []string {
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
