package client

import (
	"cashout/internal/ai"
	"cashout/internal/db"
	"cashout/internal/repository"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	// Dev Purpose, telegram usernames
	AuthEnabled  bool
	AllowedUsers map[string]struct{}
}

type Client struct {
	Logger       *logrus.Logger
	Repositories                 Repositories
	LLM                          ai.LLM
	Config                       Config
	CommandCounter               metric.Int64Counter
	TransactionOperationsCounter metric.Int64Counter
	TransactionOperationDuration metric.Float64Histogram
}

type Repositories struct {
	Users        repository.Users
	Transactions repository.Transactions
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

	// Initialize OpenTelemetry Meter and Counter
	meter := otel.Meter("cashout/client")
	commandCounter, err := meter.Int64Counter(
		"command.invocations.total",
		metric.WithDescription("Counts the number of times each bot command is invoked."),
		metric.WithUnit("{invocations}"),
	)
	if err != nil {
		// This is a programming error if initialization fails.
		// Depending on policy, could panic or log fatally.
		// For now, log fatally as the application might not behave as expected without metrics.
		logger.Fatalf("Failed to create command invocations counter: %v", err)
	}

	transactionOpsCounter, err := meter.Int64Counter(
		"transaction.operations.total",
		metric.WithDescription("Counts the number of transaction operations."),
		metric.WithUnit("{operations}"),
	)
	if err != nil {
		logger.Fatalf("Failed to create transaction operations counter: %v", err)
	}

	transactionOpDuration, err := meter.Float64Histogram(
		"transaction.operation.duration.seconds",
		metric.WithDescription("Measures the duration of transaction operations in seconds."),
		metric.WithUnit("s"),
	)
	if err != nil {
		logger.Fatalf("Failed to create transaction operation duration histogram: %v", err)
	}

	return &Client{
		Logger: logger,
		Config: config,
		Repositories: Repositories{
			Users:        repository.Users{Repository: repo},
			Transactions: repository.Transactions{Repository: repo},
		},
		LLM:                          llm,
		CommandCounter:               commandCounter,
		TransactionOperationsCounter: transactionOpsCounter,
		TransactionOperationDuration: transactionOpDuration,
	}
}
