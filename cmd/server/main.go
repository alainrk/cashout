package main

import (
	"cashout/internal/ai"
	"cashout/internal/client"
	"cashout/internal/db"
	"cashout/internal/logging"
	"context"
	"fmt"
	stdlog "log" // Renamed to avoid conflict with otel/log
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// initOtel initializes OpenTelemetry logging and metrics.
func initOtel(ctx context.Context, serviceName string, logger *logrus.Logger) (shutdown func(context.Context) error, err error) {
	// Resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("0.1.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenTelemetry resource: %w", err)
	}

	// Log Exporter & Provider
	logsEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT")
	if logsEndpoint == "" {
		logsEndpoint = "http://localhost:4318/v1/logs"
	}
	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(logsEndpoint),
		otlploghttp.WithInsecure(), // Assuming local collector or non-TLS setup
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}
	logger.Infof("OTLP Log exporter configured for endpoint: %s", logsEndpoint)

	logProcessor := sdklog.NewBatchProcessor(logExporter)
	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(logProcessor),
	)
	otel.SetLoggerProvider(logProvider) // Set as global logger provider
	logger.Info("OTLP Log provider configured and set globally.")

	// Metrics Exporter & Provider
	metricsEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT")
	if metricsEndpoint == "" {
		metricsEndpoint = "http://localhost:4318/v1/metrics"
	}
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(metricsEndpoint),
		otlpmetrichttp.WithInsecure(), // Assuming local collector or non-TLS setup
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}
	logger.Infof("OTLP Metric exporter configured for endpoint: %s", metricsEndpoint)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(meterProvider) // Set as global meter provider
	logger.Info("OTLP Meter provider configured and set globally.")

	// Shutdown function
	shutdown = func(ctx context.Context) error {
		var shutdownErr error
		logger.Info("Starting OpenTelemetry shutdown...")

		if err := logProvider.Shutdown(ctx); err != nil {
			currentErr := fmt.Errorf("failed to shutdown OTLP log provider: %w", err)
			logger.Errorf("%v", currentErr)
			if shutdownErr == nil {
				shutdownErr = currentErr
			} else {
				shutdownErr = fmt.Errorf("%v; %w", shutdownErr, currentErr)
			}
		}

		if err := meterProvider.Shutdown(ctx); err != nil {
			currentErr := fmt.Errorf("failed to shutdown OTLP meter provider: %w", err)
			logger.Errorf("%v", currentErr)
			if shutdownErr == nil {
				shutdownErr = currentErr
			} else {
				shutdownErr = fmt.Errorf("%v; %w", shutdownErr, currentErr)
			}
		}

		if shutdownErr != nil {
			logger.Info("OpenTelemetry shutdown completed with errors.")
		} else {
			logger.Info("OpenTelemetry shutdown completed successfully.")
		}
		return shutdownErr
	}

	logger.Info("OpenTelemetry initialized successfully.")
	return shutdown, nil
}

// This bot demonstrates some example interactions with commands ontelegram.
// It has a basic start command with a bot intro.
// It also has a source command, which sends the bot sourcecode, as a file.
func main() {
	err := godotenv.Load()
	if err != nil {
		stdlog.Fatal("Error loading .env file") // Use stdlog before custom logger is up
	}

	// Initialize logger early for use in Otel setup if needed
	logger := logging.GetLogger(os.Getenv("LOG_LEVEL"))

	// Add OTelLogrus hook
	logger.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	)))
	logger.Info("OTelLogrus hook added to the logger.")

	// Initialize OpenTelemetry
	otelServiceName := os.Getenv("OTEL_SERVICE_NAME")
	if otelServiceName == "" {
		otelServiceName = "cashout-telegram-bot"
	}
	otelShutdown, err := initOtel(context.Background(), otelServiceName, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer otelShutdown(context.Background()) // Ensure graceful shutdown

	// Get token from the environment variable
	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if token == "" {
		logger.Fatalln("TELEGRAM_BOT_API_TOKEN environment variable is empty")
	}

	// API key and endpoint
	aiApiKey := os.Getenv("DEEPSEEK_API_KEY")
	aiEndpoint := "https://api.deepseek.com/v1/chat/completions"
	llm := ai.LLM{
		Logger:   logger,
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
			logger.Fatalf("failed to start webhook: %s\n", err.Error()) // Changed panic to logger.Fatalf
		}

		err = updater.SetAllBotWebhooks(webhookDomain, &gotgbot.SetWebhookOpts{
			MaxConnections:     100,
			DropPendingUpdates: true,
			SecretToken:        webhookOpts.SecretToken,
		})
		if err != nil {
			logger.Fatalf("failed to set webhook: %s\n", err.Error()) // Changed panic to logger.Fatalf
		}
	default:
		logger.Fatalf("unknown run mode: %s\n", runMode)
	}

	logger.Infof("%s has been started in %s mode...\n", b.Username, runMode)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
