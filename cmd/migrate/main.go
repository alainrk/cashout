package main

import (
	"cashout/internal/migrations"
	"flag"
	"fmt"
	"log"
	"os"

	_ "cashout/internal/migrations/versions" // Import all migrations

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Add a new flag for environment file
	var (
		command string
		envFile string
	)

	flag.StringVar(&command, "command", "up", "Migration command (up, down, status)")
	flag.StringVar(&envFile, "env", ".env", "Environment file to load (.env, .prod.env, etc)")
	flag.Parse()

	// Load the specified environment file
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}

	// Initialize database
	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		panic("DATABASE_URL environment variable is empty")
	}

	// Connect to the database directly with GORM
	conn, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a migrator
	migrator := migrations.NewMigrator(conn)

	// Execute the requested command
	switch command {
	case "up":
		if err := migrator.MigrateUp(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("All migrations applied successfully!")
	case "down":
		if err := migrator.MigrateDown(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "status":
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
