package main

import (
	"cashout/internal/ai"
	"cashout/internal/client"
	"cashout/internal/db"
	"cashout/internal/logging"
	"log"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
)

// Create a matcher which only matches text which is not a command.
func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}

func confirmCommand(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.Trim(msg.Text, " ") == "Confirm"
}

// This bot demonstrates some example interactions with commands ontelegram.
// It has a basic start command with a bot intro.
// It also has a source command, which sends the bot sourcecode, as a file.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := logging.GetLogger(os.Getenv("LOG_LEVEL"))

	// Get token from the environment variable
	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if token == "" {
		logger.Fatalln("TELEGRAM_BOT_API_TOKEN environment variable is empty")
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
		logger.Fatalln("DATABASE_URL environment variable is empty")
	}

	db, err := db.NewDB(postgresURL)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %s\n", err.Error())
	}

	defer db.Close()

	// Initialize client
	c := client.NewClient(logger, db, llm)

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		logger.Fatalf("failed to create new bot: %s\n", err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Errorf("an error occurred while handling update: %s\n", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	client.SetupHandlers(dispatcher, c)

	///////////////////////////////////////

	webhookDomain := os.Getenv("WEBHOOK_DOMAIN")
	if webhookDomain == "" {
		log.Fatalln("WEBHOOK_DOMAIN environment variable is empty")
	}

	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		panic("WEBHOOK_SECRET environment variable is empty")
	}

	// Start the webhook server, but before start the server so we're ready when Telegram starts sending updates.
	webhookOpts := ext.WebhookOpts{
		ListenAddr:  "localhost:8080", // TODO: Put it into config
		SecretToken: webhookSecret,
	}

	// The bot's urlPath can be anything.
	// It's a good idea to contain the bot token, as that makes it very difficult for outside
	// parties to find the update endpoint (which would allow them to inject their own updates).
	err = updater.StartWebhook(b, "cashout/"+token, webhookOpts)
	if err != nil {
		panic("failed to start webhook: " + err.Error())
	}

	err = updater.SetAllBotWebhooks(webhookDomain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookOpts.SecretToken,
	})
	if err != nil {
		panic("failed to set webhook: " + err.Error())
	}

	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
