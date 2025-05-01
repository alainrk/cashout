package main

import (
	"happypoor/internal/ai"
	"happypoor/internal/client"
	"happypoor/internal/db"
	"log"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
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

	// Get token from the environment variable
	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if token == "" {
		panic("TELEGRAM_BOT_API_TOKEN environment variable is empty")
	}

	webhookDomain := os.Getenv("WEBHOOK_DOMAIN")
	if webhookDomain == "" {
		panic("WEBHOOK_DOMAIN environment variable is empty")
	}

	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		panic("WEBHOOK_SECRET environment variable is empty")
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

	// Create updater and dispatcher for webhook management
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	// // Create updater and dispatcher.
	// dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
	// 	// If an error is returned by a handler, log it and continue going.
	// 	Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
	// 		log.Println("an error occurred while handling update:", err.Error())
	// 		return ext.DispatcherActionNoop
	// 	},
	// 	MaxRoutines: ext.DefaultMaxRoutines,
	// })
	//
	// updater := ext.NewUpdater(dispatcher, nil)

	// Top-level message for LLM goes into AddTransaction and gets the expense/income intent from user session state.
	dispatcher.AddHandler(handlers.NewMessage(noCommands, c.FreeTextRouter))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("transactions.new."), c.AddTransactionIntent))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("transactions.edit."), c.EditTransactionIntent))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.cancel"), c.Cancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.confirm"), c.Confirm))

	dispatcher.AddHandler(handlers.NewCommand("list", c.ListTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.year."), c.ListYearNavigation))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.month."), c.ListMonthTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.page."), c.ListTransactionPage))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.page."), c.DeleteTransactionPage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.confirm."), c.DeleteTransactionConfirm))

	dispatcher.AddHandler(handlers.NewCommand("cancel", c.Cancel))
	dispatcher.AddHandler(handlers.NewCommand("delete", c.DeleteTransactions))
	dispatcher.AddHandler(handlers.NewCommand("start", c.Start))
	dispatcher.AddHandler(handlers.NewCommand("new", c.Start))
	dispatcher.AddHandler(handlers.NewCommand("month", c.MonthRecap))
	dispatcher.AddHandler(handlers.NewCommand("year", c.YearRecap))

	// TODO:
	// dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.month"), c.MonthRecap))
	// dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.year"), c.YearRecap))
	// dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.list"), c.ListTransactions))
	// dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.delete"), c.DeleteTransactions))

	// Start the webhook server, but before start the server so we're ready when Telegram starts sending updates.
	webhookOpts := ext.WebhookOpts{
		ListenAddr:  "localhost:3666", // TODo: Put it into config
		SecretToken: webhookSecret,
	}

	// The bot's urlPath can be anything.
	// It's a good idea to contain the bot token, as that makes it very difficult for outside
	// parties to find the update endpoint (which would allow them to inject their own updates).
	err = updater.StartWebhook(b, "happywebhook/"+token, webhookOpts)
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

	// Start receiving updates.
	// err = updater.StartPolling(b, &ext.PollingOpts{
	// 	DropPendingUpdates: true,
	// 	GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
	// 		Timeout: 9,
	// 		RequestOpts: &gotgbot.RequestOpts{
	// 			Timeout: time.Second * 10,
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	panic("failed to start polling: " + err.Error())
	// }
	//
	// log.Printf("%s has been started...\n", b.Username)
	//
	// // Idle, to keep updates coming in, and avoid bot stopping.
	// updater.Idle()
}
