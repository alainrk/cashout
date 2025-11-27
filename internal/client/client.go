package client

import (
	"os"
	"strings"

	"cashout/internal/ai"
	"cashout/internal/db"
	"cashout/internal/repository"

	"github.com/sirupsen/logrus"
)

const MinYearAllowed = 2015

type Config struct {
	// Dev Purpose, telegram usernames
	AuthEnabled     bool
	AllowedUsers    map[string]struct{}
	WebDashboardURL string
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
		for u := range strings.SplitSeq(usernames, ",") {
			config.AllowedUsers[u] = struct{}{}
		}
	}

	config.WebDashboardURL = os.Getenv("WEB_DASHBOARD_URL")

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
