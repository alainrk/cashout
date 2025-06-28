package main

import (
	"cashout/internal/db"
	"cashout/internal/model"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type Seeder struct {
	db       *db.DB
	userTgID int64
}

func NewSeeder(database *db.DB, userTgID int64) *Seeder {
	return &Seeder{
		db:       database,
		userTgID: userTgID,
	}
}

func (s *Seeder) SeedTransactions() error {
	// First, check if user exists
	user, err := s.db.GetUser(s.userTgID)
	if err != nil {
		return fmt.Errorf("user with TG ID %d not found: %w", s.userTgID, err)
	}
	fmt.Printf("Found user: %s (TG ID: %d)\n", user.Name, user.TgID)

	// Delete all existing transactions for this user (idempotent)
	fmt.Println("Deleting existing transactions...")
	if err := s.deleteUserTransactions(); err != nil {
		return fmt.Errorf("failed to delete existing transactions: %w", err)
	}

	// Seed the faker
	gofakeit.Seed(time.Now().UnixNano())

	// Calculate date range (5 years ago to today)
	endDate := time.Now()
	startDate := endDate.AddDate(-5, 0, 0)

	fmt.Printf("Generating transactions from %s to %s\n",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	// Generate transactions
	transactions := s.generateTransactions(startDate, endDate)

	// Insert transactions in batches
	fmt.Printf("Inserting %d transactions...\n", len(transactions))
	batchSize := 100
	for i := 0; i < len(transactions); i += batchSize {
		end := i + batchSize
		if end > len(transactions) {
			end = len(transactions)
		}

		batch := transactions[i:end]
		for _, t := range batch {
			if err := s.db.CreateTransaction(&t); err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}
		}

		fmt.Printf("Progress: %d/%d transactions inserted\n", end, len(transactions))
	}

	return nil
}

func (s *Seeder) deleteUserTransactions() error {
	// Get all transactions for the user and delete them
	transactions, err := s.db.GetUserTransactions(s.userTgID)
	if err != nil {
		return err
	}

	for _, t := range transactions {
		if err := s.db.DeleteTransaction(t.ID); err != nil {
			return fmt.Errorf("failed to delete transaction %d: %w", t.ID, err)
		}
	}

	fmt.Printf("Deleted %d existing transactions\n", len(transactions))
	return nil
}

func (s *Seeder) generateTransactions(startDate, endDate time.Time) []model.Transaction {
	var transactions []model.Transaction

	// All expense categories
	expenseCategories := []model.TransactionCategory{
		model.CategoryCar,
		model.CategoryClothes,
		model.CategoryGrocery,
		model.CategoryHouse,
		model.CategoryBills,
		model.CategoryEntertainment,
		model.CategorySport,
		model.CategoryEatingOut,
		model.CategoryTransport,
		model.CategoryLearning,
		model.CategoryToiletry,
		model.CategoryHealth,
		model.CategoryTech,
		model.CategoryGifts,
		model.CategoryTravel,
		model.CategoryPets,
		model.CategoryOtherExpenses,
	}

	// Income categories
	incomeCategories := []model.TransactionCategory{
		model.CategorySalary,
		model.CategoryOtherIncomes,
	}

	// Generate transactions for each day in the range
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		// Random number of transactions per day (0-3)
		numTransactions := rand.Intn(4)

		for i := 0; i < numTransactions; i++ {
			// 90% chance of expense, 10% chance of income
			isExpense := rand.Float32() < 0.9

			var transaction model.Transaction
			transaction.TgID = s.userTgID
			transaction.Date = d
			transaction.Currency = model.CurrencyEUR

			// In generateTransactions function...
			if isExpense {
				transaction.Type = model.TypeExpense
				transaction.Category = expenseCategories[rand.Intn(len(expenseCategories))]
				transaction.Amount = s.generateExpenseAmount(transaction.Category)
				transaction.Description = s.generateExpenseDescription(transaction.Category)
			} else {
				transaction.Type = model.TypeIncome
				transaction.Category = incomeCategories[rand.Intn(len(incomeCategories))]
				transaction.Amount = s.generateIncomeAmount(transaction.Category)
				// Pass the transaction date to generate a correct description
				transaction.Description = s.generateIncomeDescription(transaction.Category, transaction.Date) // Changed line
			}

			transactions = append(transactions, transaction)
		}
	}

	// Ensure at least one salary per month
	s.ensureMonthlySalaries(&transactions, startDate, endDate)

	// Shuffle transactions to make them more realistic
	rand.Shuffle(len(transactions), func(i, j int) {
		transactions[i], transactions[j] = transactions[j], transactions[i]
	})

	return transactions
}

func (s *Seeder) generateExpenseAmount(category model.TransactionCategory) float64 {
	// Generate realistic amounts based on category
	switch category {
	case model.CategoryGrocery:
		return gofakeit.Float64Range(10, 150)
	case model.CategoryCar:
		return gofakeit.Float64Range(20, 200)
	case model.CategoryClothes:
		return gofakeit.Float64Range(15, 300)
	case model.CategoryHouse:
		return gofakeit.Float64Range(50, 2000)
	case model.CategoryBills:
		return gofakeit.Float64Range(30, 300)
	case model.CategoryEntertainment:
		return gofakeit.Float64Range(10, 100)
	case model.CategorySport:
		return gofakeit.Float64Range(20, 150)
	case model.CategoryEatingOut:
		return gofakeit.Float64Range(15, 80)
	case model.CategoryTransport:
		return gofakeit.Float64Range(2, 50)
	case model.CategoryLearning:
		return gofakeit.Float64Range(20, 500)
	case model.CategoryToiletry:
		return gofakeit.Float64Range(5, 50)
	case model.CategoryHealth:
		return gofakeit.Float64Range(20, 200)
	case model.CategoryTech:
		return gofakeit.Float64Range(30, 1000)
	case model.CategoryGifts:
		return gofakeit.Float64Range(20, 200)
	case model.CategoryTravel:
		return gofakeit.Float64Range(50, 1500)
	default:
		return gofakeit.Float64Range(10, 100)
	}
}

func (s *Seeder) generateIncomeAmount(category model.TransactionCategory) float64 {
	switch category {
	case model.CategorySalary:
		return gofakeit.Float64Range(2000, 5000)
	case model.CategoryOtherIncomes:
		return gofakeit.Float64Range(50, 500)
	default:
		return gofakeit.Float64Range(100, 1000)
	}
}

func (s *Seeder) generateExpenseDescription(category model.TransactionCategory) string {
	switch category {
	case model.CategoryGrocery:
		stores := []string{"Supermarket", "Local Market", "Grocery Store", "Farmers Market"}
		return gofakeit.RandomString(stores)
	case model.CategoryCar:
		return gofakeit.RandomString([]string{"Gas", "Car wash", "Parking", "Maintenance", "Insurance"})
	case model.CategoryClothes:
		return gofakeit.Company() + " Store"
	case model.CategoryHouse:
		return gofakeit.RandomString([]string{"Rent", "Utilities", "Internet", "Cleaning supplies", "Furniture"})
	case model.CategoryBills:
		return gofakeit.RandomString([]string{"Electricity", "Water", "Gas", "Phone", "Internet"})
	case model.CategoryEntertainment:
		return gofakeit.RandomString([]string{"Cinema", "Concert", "Theatre", "Museum", "Netflix"})
	case model.CategorySport:
		return gofakeit.RandomString([]string{"Gym membership", "Swimming pool", "Sports equipment", "Fitness class"})
	case model.CategoryEatingOut:
		return gofakeit.Company() + " Restaurant"
	case model.CategoryTransport:
		return gofakeit.RandomString([]string{"Bus ticket", "Metro", "Taxi", "Uber", "Train"})
	case model.CategoryLearning:
		return gofakeit.RandomString([]string{"Online course", "Books", "Workshop", "Tutorial", "Conference"})
	case model.CategoryToiletry:
		return gofakeit.RandomString([]string{"Shampoo", "Toothpaste", "Soap", "Cosmetics", "Toiletries"})
	case model.CategoryHealth:
		return gofakeit.RandomString([]string{"Doctor visit", "Pharmacy", "Vitamins", "Medical check-up"})
	case model.CategoryTech:
		return gofakeit.RandomString([]string{"Software", "Hardware", "Gadget", "Phone accessory", "Computer parts"})
	case model.CategoryGifts:
		return "Gift for " + gofakeit.FirstName()
	case model.CategoryTravel:
		return "Trip to " + gofakeit.City()
	case model.CategoryPets:
		descriptions := []string{"Pet food", "Vet visit", "Pet toys", "Grooming", "Medication", "Treats"}
		return gofakeit.RandomString(descriptions)
	case model.CategoryOtherExpenses:
		descriptions := []string{"Miscellaneous goods", "Service payment", "Online purchase", "General expense"}
		return gofakeit.RandomString(descriptions)
	default:
		// Use a more descriptive fallback than a single word
		return gofakeit.ProductName()
	}
}

func (s *Seeder) generateIncomeDescription(category model.TransactionCategory, date time.Time) string { // Signature updated
	switch category {
	case model.CategorySalary:
		// Use the actual month from the transaction's date
		return fmt.Sprintf("%s salary", date.Format("January"))
	case model.CategoryOtherIncomes:
		return gofakeit.RandomString([]string{"Freelance work", "Bonus", "Refund", "Gift received", "Investment return"})
	default:
		return "Income"
	}
}

func (s *Seeder) ensureMonthlySalaries(transactions *[]model.Transaction, startDate, endDate time.Time) {
	// Create a map to track which months have salaries
	salaryMonths := make(map[string]bool)

	// Check existing transactions for salaries
	for _, t := range *transactions {
		if t.Type == model.TypeIncome && t.Category == model.CategorySalary {
			monthKey := fmt.Sprintf("%d-%02d", t.Date.Year(), t.Date.Month())
			salaryMonths[monthKey] = true
		}
	}

	// Add missing salaries
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
		monthKey := fmt.Sprintf("%d-%02d", d.Year(), d.Month())

		if !salaryMonths[monthKey] {
			// Add salary on a random day between 25-28 of each month
			salaryDay := 25 + rand.Intn(4)
			salaryDate := time.Date(d.Year(), d.Month(), salaryDay, 0, 0, 0, 0, d.Location())

			// Make sure the salary date is within our range
			if salaryDate.After(endDate) {
				salaryDate = endDate
			}
			if salaryDate.Before(startDate) {
				continue
			}

			salary := model.Transaction{
				TgID:        s.userTgID,
				Date:        salaryDate,
				Type:        model.TypeIncome,
				Category:    model.CategorySalary,
				Amount:      s.generateIncomeAmount(model.CategorySalary),
				Currency:    model.CurrencyEUR,
				Description: fmt.Sprintf("%s salary", salaryDate.Format("January")),
			}

			*transactions = append(*transactions, salary)
		}
	}
}
