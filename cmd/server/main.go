package main

import (
	"happypoor/client"
	"happypoor/internal/db"
	"log"
	"os"
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
	c := client.Client{Db: db, Store: &client.Store{Db: db}}

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
	dispatcher.AddHandler(handlers.NewCommand("cancel", c.Cancel))
	dispatcher.AddHandler(handlers.NewMessage(noCommands, c.Message))

	// dispatcher.AddHandler(handlers.NewConversation(
	// 	[]ext.Handler{handlers.NewCommand("start", c.Start)},
	// 	map[string][]ext.Handler{
	// 		string(model.StateNormal): {handlers.NewMessage(noCommands, c.Message)},
	// 	},
	// 	&handlers.ConversationOpts{
	// 		Exits: []ext.Handler{handlers.NewCommand("cancel", c.Cancel)},
	// 		// StateStorage: c.Store,
	// 		StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
	// 		Fallbacks:    []ext.Handler{handlers.NewMessage(noCommands, c.Message)},
	// 		AllowReEntry: true,
	// 	},
	// ))

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
