package main

import (
	"happypoor/internal/ai"
	"happypoor/internal/client"
	"happypoor/internal/db"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
)

// Create a matcher which only matches text which is not a command.
func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}

func cancelCommand(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.Trim(msg.Text, " ") == "Cancel"
}

func addIncome(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.Trim(msg.Text, " ") == "Add Income"
}

func addExpense(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.Trim(msg.Text, " ") == "Add Expense"
}

func confirmCommand(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.Trim(msg.Text, " ") == "Confirm"
}

func amendCommand(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.Trim(msg.Text, " ") == "Edit"
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
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewCommand("start", c.Start))
	dispatcher.AddHandler(handlers.NewCommand("income", c.AddIncomeIntent))
	dispatcher.AddHandler(handlers.NewCommand("expense", c.AddExpenseIntent))
	dispatcher.AddHandler(handlers.NewCommand("cancel", c.Cancel))
	dispatcher.AddHandler(handlers.NewCommand("month", c.MonthRecap))
	dispatcher.AddHandler(handlers.NewMessage(cancelCommand, c.Cancel))
	dispatcher.AddHandler(handlers.NewMessage(confirmCommand, c.Confirm))
	dispatcher.AddHandler(handlers.NewMessage(amendCommand, c.AmendTransaction))
	dispatcher.AddHandler(handlers.NewMessage(addIncome, c.AddIncomeIntent))
	dispatcher.AddHandler(handlers.NewMessage(addExpense, c.AddExpenseIntent))
	dispatcher.AddHandler(handlers.NewMessage(noCommands, c.AddTransaction))

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
