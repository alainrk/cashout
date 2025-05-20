package main

import (
	"context"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestInitOtelInitializesAndShutsDownCleanly(t *testing.T) {
	// Set environment variables for testing
	// Using localhost with a port that's unlikely to be in use and non-standard for OTLP
	// to prevent actual data export and ensure the URLs are parseable.
	originalLogsEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT")
	originalMetricsEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT")
	originalServiceName := os.Getenv("OTEL_SERVICE_NAME")

	t.Setenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT", "http://localhost:1/v1/logs")
	t.Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "http://localhost:1/v1/metrics")
	t.Setenv("OTEL_SERVICE_NAME", "test-cashout-service")

	// Restore original environment variables after the test
	defer func() {
		os.Setenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT", originalLogsEndpoint)
		os.Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", originalMetricsEndpoint)
		os.Setenv("OTEL_SERVICE_NAME", originalServiceName)
	}()

	// Create a minimal logger for initOtel to use
	testLogger := logrus.New()
	testLogger.SetOutput(os.Stderr) // Or io.Discard if logs are too noisy for tests
	testLogger.SetLevel(logrus.DebugLevel)

	// Call initOtel
	otelShutdown, err := initOtel(context.Background(), "test-service", testLogger)

	// Assert that initOtel completed without error
	require.NoError(t, err, "initOtel should not return an error")
	require.NotNil(t, otelShutdown, "initOtel should return a non-nil shutdown function")

	// Check that global providers are set
	// Note: These are global states, which can be problematic for parallel tests.
	// For this test, we assume it's acceptable or tests are run sequentially.
	loggerProvider := otel.GetLoggerProvider()
	meterProvider := otel.GetMeterProvider()

	assert.NotNil(t, loggerProvider, "Global LoggerProvider should be set by initOtel")
	assert.NotNil(t, meterProvider, "Global MeterProvider should be set by initOtel")

	// Call the shutdown function and assert it returns no error
	errShutdown := otelShutdown(context.Background())
	assert.NoError(t, errShutdown, "otelShutdown should not return an error")

	// Additionally, you could try to reset global providers after the test,
	// though this is often tricky and not perfectly clean.
	// For example:
	// otel.SetLoggerProvider(otel.NewNoopLoggerProvider())
	// otel.SetMeterProvider(otel.NewNoopMeterProvider())
	// This helps isolate tests if they are run in the same process, but `initOtel`
	// itself uses otel.SetLoggerProvider and otel.SetMeterProvider, so subsequent
	// calls in other tests would override this. The -p 1 flag for go test is a more
	// robust way to handle global state issues if they arise.
}
