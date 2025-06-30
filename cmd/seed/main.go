package main

import (
	"cashout/internal/db"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	var envFile string
	flag.StringVar(&envFile, "env", ".env", "Environment file to load (.env, .prod.env, etc)")
	flag.Parse()

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}

	userTgIDStr := os.Getenv("SEED_USER_TG_ID")
	if userTgIDStr == "" {
		log.Fatal("SEED_USER_TG_ID environment variable is not set. Please set it to the Telegram ID of the user you want to seed transactions for.")
	}

	userTgID, err := strconv.ParseInt(userTgIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid SEED_USER_TG_ID: %v", err)
	}

	// Initialize database
	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		log.Fatal("DATABASE_URL environment variable is empty")
	}

	database, err := db.NewDB(postgresURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		err = errors.Join(err, database.Close())
	}()

	// Create and run the seeder
	seeder := NewSeeder(database, userTgID)

	fmt.Printf("Starting transaction seed for user with TG ID: %d\n", userTgID)

	if err := seeder.SeedTransactions(); err != nil {
		log.Fatalf("Failed to seed transactions: %v", err)
	}

	fmt.Println("Transaction seeding completed successfully!")
}
