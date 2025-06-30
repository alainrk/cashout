package main

import (
	"cashout/internal/ai"
	"cashout/internal/client"
	"cashout/internal/db"
	"cashout/internal/logging"
	"cashout/internal/scheduler"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joho/godotenv"

	_ "github.com/go-co-op/gocron" // Add to go.mod
)

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

	// OpenAI API Compatible LLM Setup
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

	runMode := strings.ToLower(os.Getenv("RUN_MODE"))

	switch runMode {
	case "polling":
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
			logger.Fatalf("failed to start polling: %s\n", err.Error())
		}
	case "webhook":
		webhookDomain := os.Getenv("WEBHOOK_DOMAIN")
		if webhookDomain == "" {
			logger.Fatalln("WEBHOOK_DOMAIN environment variable is empty")
		}

		webhookSecret := os.Getenv("WEBHOOK_SECRET")
		if webhookSecret == "" {
			logger.Fatalln("WEBHOOK_SECRET environment variable is empty")
		}

		webhookHost := os.Getenv("WEBHOOK_HOST")
		if webhookHost == "" {
			webhookHost = "0.0.0.0"
		}

		webhookPort := os.Getenv("WEBHOOK_PORT")
		if webhookPort == "" {
			webhookPort = "8080"
		}

		// Start the webhook server, but before start the server so we're ready when Telegram starts sending updates.
		webhookOpts := ext.WebhookOpts{
			ListenAddr:  webhookHost + ":" + webhookPort,
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
	default:
		logger.Fatalf("unknown run mode: %s\n", runMode)
	}

	logger.Infof("%s has been started in %s mode...\n", b.Username, runMode)

	// Initialize scheduler for automated reminders
	sched := scheduler.NewScheduler(b, c.Repositories, logger)
	sched.Start()
	defer sched.Stop()

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
