package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/metric" // Not directly used for Add/Record if instrument is typed
)

// TestAIMetricsInitializationAndInteraction tests that the global AI metric instruments
// are initialized (non-nil) and can be interacted with.
func TestAIMetricsInitializationAndInteraction(t *testing.T) {
	// 1. Test Metric Instrument Initialization
	// These are global variables initialized in the package's init() function.
	// We expect them to be non-nil if the init() function ran correctly.
	assert.NotNil(t, aiAPICallsCounter, "aiAPICallsCounter should be initialized")
	assert.NotNil(t, aiAPICallDuration, "aiAPICallDuration should be initialized")

	// 2. Test Basic Metric Interaction
	// Ensure Add/Record calls don't panic. This relies on the global MeterProvider
	// being a no-op provider if not fully configured, which is default OTel SDK behavior.

	// Test aiAPICallsCounter
	// We need to check if the counter is nil before using it, as the init() function
	// might have failed to initialize it (e.g., if otel.Meter() had an issue, though unlikely for Noop).
	if aiAPICallsCounter != nil {
		assert.NotPanics(t, func() {
			aiAPICallsCounter.Add(context.Background(), 1, attribute.String("status", "test_success"))
		}, "aiAPICallsCounter.Add should not panic")
	} else {
		t.Log("aiAPICallsCounter is nil, skipping interaction test. This might indicate an issue in the init() function or test setup.")
	}

	// Test aiAPICallDuration
	if aiAPICallDuration != nil {
		assert.NotPanics(t, func() {
			aiAPICallDuration.Record(context.Background(), 0.5, attribute.String("status", "test_success"))
		}, "aiAPICallDuration.Record should not panic")
	} else {
		t.Log("aiAPICallDuration is nil, skipping interaction test. This might indicate an issue in the init() function or test setup.")
	}
}
