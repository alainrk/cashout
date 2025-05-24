package client

import (
	"cashout/internal/ai"
	"cashout/internal/db"
	"cashout/internal/repository"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	// Dev Purpose, telegram usernames
	AuthEnabled  bool
	AllowedUsers map[string]struct{}
}

type Client struct {
	Logger       *logrus.Logger
	Repositories Repositories
	LLM          ai.LLM
	Config       Config
}

type Repositories struct {
	Users        repository.Users
	Transactions repository.Transactions
	Reminders    repository.Reminders
}

func NewClient(logger *logrus.Logger, db *db.DB, llm ai.LLM) *Client {
	config := Config{
		AllowedUsers: make(map[string]struct{}),
	}

	usernames := os.Getenv("ALLOWED_USERS")
	if usernames != "" {
		config.AuthEnabled = true
		for _, u := range strings.Split(usernames, ",") {
			config.AllowedUsers[u] = struct{}{}
		}
	}

	// For repositories structs embedding common fields
	repo := repository.Repository{
		DB:     db,
		Logger: logger,
	}

	return &Client{
		Logger: logger,
		Config: config,
		Repositories: Repositories{
			Users:        repository.Users{Repository: repo},
			Transactions: repository.Transactions{Repository: repo},
			Reminders:    repository.Reminders{Repository: repo},
		},
		LLM: llm,
	}
}
