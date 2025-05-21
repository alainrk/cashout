package repository

import (
	"cashout/internal/db"
	"log" // Standard log for init errors

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	dbOperationsCounter metric.Int64Counter
	dbOperationDuration metric.Float64Histogram
)

func init() {
	meter := otel.Meter("cashout/repository")
	var err error
	dbOperationsCounter, err = meter.Int64Counter(
		"db.operations.total",
		metric.WithDescription("Counts the number of database operations."),
	)
	if err != nil {
		log.Printf("Error initializing dbOperationsCounter: %v\n", err) // Or panic
	}

	dbOperationDuration, err = meter.Float64Histogram(
		"db.operation.duration.seconds",
		metric.WithDescription("Measures the duration of database operations in seconds."),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Printf("Error initializing dbOperationDuration: %v\n", err) // Or panic
	}
}

type Repository struct {
	DB     *db.DB
	Logger *logrus.Logger
}
