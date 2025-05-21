package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/metric" // Not directly used for Add/Record if instrument is typed
)

// TestRepositoryMetricsInitializationAndInteraction tests that the global repository metric instruments
// are initialized (non-nil) and can be interacted with.
func TestRepositoryMetricsInitializationAndInteraction(t *testing.T) {
	// 1. Test Metric Instrument Initialization
	// These are global variables initialized in the package's init() function.
	assert.NotNil(t, dbOperationsCounter, "dbOperationsCounter should be initialized")
	assert.NotNil(t, dbOperationDuration, "dbOperationDuration should be initialized")

	// 2. Test Basic Metric Interaction
	// Ensure Add/Record calls don't panic. This relies on the global MeterProvider
	// being a no-op provider if not fully configured, which is default OTel SDK behavior.

	// Test dbOperationsCounter
	if dbOperationsCounter != nil {
		assert.NotPanics(t, func() {
			dbOperationsCounter.Add(context.Background(), 1,
				attribute.String("db.table", "test_table"),
				attribute.String("db.operation", "test_op"),
				attribute.String("status", "success"),
			)
		}, "dbOperationsCounter.Add should not panic")
	} else {
		t.Log("dbOperationsCounter is nil, skipping interaction test. This might indicate an issue in the init() function or test setup.")
	}

	// Test dbOperationDuration
	if dbOperationDuration != nil {
		assert.NotPanics(t, func() {
			dbOperationDuration.Record(context.Background(), 0.05,
				attribute.String("db.table", "test_table"),
				attribute.String("db.operation", "test_op"),
				attribute.String("status", "success"),
			)
		}, "dbOperationDuration.Record should not panic")
	} else {
		t.Log("dbOperationDuration is nil, skipping interaction test. This might indicate an issue in the init() function or test setup.")
	}
}
