package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cashout/internal/ai"
	"cashout/internal/db"
	"cashout/internal/email"
	"cashout/internal/logging"
	"cashout/internal/repository"
	"cashout/internal/web"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := logging.GetLogger(os.Getenv("LOG_LEVEL"))

	// Initialize email service
	emailService, err := email.NewEmailService(os.Getenv("BREVO_API_KEY"), os.Getenv("EMAIL_FROM_NAME"), os.Getenv("EMAIL_FROM_ADDRESS"))
	if err != nil {
		logger.Fatalf("Failed to initialize email service: %s\n", err.Error())
	}

	// Get token from the environment variable
	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if token == "" {
		logger.Fatalln("TELEGRAM_BOT_API_TOKEN environment variable is empty")
	}

	// Create bot for sending auth codes
	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		logger.Fatalf("failed to create new bot: %s\n", err.Error())
	}

	// OpenAI API Compatible LLM Setup (if needed for dashboard features)
	llm := ai.LLM{
		Logger:   logger,
		APIKey:   os.Getenv("OPENAI_API_KEY"),
		Model:    os.Getenv("LLM_MODEL"),
		Endpoint: fmt.Sprintf("%s/chat/completions", os.Getenv("OPENAI_BASE_URL")),
	}

	// Initialize database
	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		logger.Fatalln("DATABASE_URL environment variable is empty")
	}

	database, err := db.NewDB(postgresURL)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %s\n", err.Error())
	}

	defer func() {
		err = errors.Join(err, database.Close())
	}()

	// For repositories structs embedding common fields
	repo := repository.Repository{
		DB:     database,
		Logger: logger,
	}

	// Initialize WebAuthn repository
	webAuthnRepo, err := repository.NewWebAuthn(repo)
	if err != nil {
		logger.Fatalf("Failed to initialize WebAuthn repository: %s\n", err.Error())
	}

	repositories := web.Repositories{
		Users:        repository.Users{Repository: repo},
		Transactions: repository.Transactions{Repository: repo},
		Auth:         repository.Auth{Repository: repo},
		WebAuthn:     webAuthnRepo,
	}

	// Start periodic WebAuthn session cleanup (every hour)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := webAuthnRepo.CleanupExpiredSessions(); err != nil {
				logger.Errorf("Failed to cleanup expired WebAuthn sessions: %v", err)
			} else {
				logger.Debugf("Successfully cleaned up expired WebAuthn sessions")
			}
		}
	}()

	// Initialize web server
	webServer := web.NewServer(logger, repositories, bot, llm, emailService)

	// Get web server configuration
	webHost := os.Getenv("WEB_HOST")
	if webHost == "" {
		webHost = "localhost"
	}

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = "8081"
	}

	addr := fmt.Sprintf("%s:%s", webHost, webPort)
	logger.Infof("Starting web server on %s", addr)

	if err := http.ListenAndServe(addr, web.Router(webServer)); err != nil {
		logger.Fatalf("Web server failed: %s", err.Error())
	}
}
