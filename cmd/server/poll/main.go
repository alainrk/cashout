package main

import (
	"happypoor/internal/ai"
	"happypoor/internal/client"
	"happypoor/internal/db"
	"happypoor/internal/handlers"
	"log"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joho/godotenv"
)

// This bot demonstrates some example interactions with commands ontelegram.
// It has a basic start command with a bot intro.
// It also has a source command, which sends the bot sourcecode, as a file.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get token from the environment variable
	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if token == "" {
		panic("TELEGRAM_BOT_API_TOKEN environment variable is empty")
	}

	// API key and endpoint
	aiApiKey := os.Getenv("DEEPSEEK_API_KEY")
	aiEndpoint := "https://api.deepseek.com/v1/chat/completions"
	llm := ai.LLM{
		APIKey:   aiApiKey,
		Endpoint: aiEndpoint,
	}

	// Initialize database
	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		panic("DATABASE_URL environment variable is empty")
	}
	db, err := db.NewDB(postgresURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize client
	c := client.NewClient(db, llm)

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	handlers.SetupHandlers(dispatcher, c)

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	log.Printf("%s has been started...\n", b.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
