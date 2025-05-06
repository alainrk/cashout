package client

import (
	"cashout/internal/ai"
	"cashout/internal/db"
	"cashout/internal/repository"
	"os"
	"strings"
)

type Config struct {
	// Dev Purpose, telegram usernames
	AuthEnabled  bool
	AllowedUsers map[string]struct{}
}

type Client struct {
	Repositories Repositories
	LLM          ai.LLM
	Config       Config
}

type Repositories struct {
	Users        repository.Users
	Transactions repository.Transactions
}

func NewClient(db *db.DB, llm ai.LLM) *Client {
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

	return &Client{
		Config: config,
		Repositories: Repositories{
			Users: repository.Users{DB: db},
			Transactions: repository.Transactions{
				DB: db,
			},
		},
		LLM: llm,
	}
}
