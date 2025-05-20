package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

func TestOtelLogrusHookIsAdded(t *testing.T) {
	logger := logrus.New()
	// Create the hook, similar to how it might be done in main.go or other setup code.
	// Using AllLevels for simplicity in the test.
	otelHook := otellogrus.NewHook(otellogrus.WithLevels(logrus.AllLevels...))
	logger.AddHook(otelHook)

	found := false
	// logrus.Logger.Hooks is a map where keys are log levels and values are lists of hooks.
	for _, levelHooks := range logger.Hooks {
		for _, hook := range levelHooks {
			if _, ok := hook.(*otellogrus.Hook); ok {
				found = true
				break // Found the hook, exit inner loop
			}
		}
		if found {
			break // Found the hook, exit outer loop
		}
	}

	assert.True(t, found, "otellogrus.Hook should be present in logger hooks")
}
