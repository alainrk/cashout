package repository

import (
	"cashout/internal/model"
	"context"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Users struct {
	Repository
}

func (r *Users) GetByUsername(username string) (user model.User, found bool, err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		// "found" could also be an attribute if desired, e.g., attribute.Bool("found", found)
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "users"),
			attribute.String("db.operation", "get_by_username"),
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	dbUser, err := r.DB.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			return model.User{}, false, nil // Not an operational error for metrics purposes, but 'found' is false.
		}
		return model.User{}, false, err // Actual error
	}
	return *dbUser, true, nil
}

func (r *Users) UpsertWithContext(botUser gotgbot.User, session model.UserSession) (err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "users"),
			attribute.String("db.operation", "upsert_with_context"), // This implies an upsert logic in DB.SetUser
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	name := botUser.FirstName
	name = strings.Trim(name, " ")
	if name == "" {
		name = botUser.Username
	}

	err = r.DB.SetUser(&model.User{
		TgID:        botUser.Id,
		Name:        name,
		Session:     session,
		TgUsername:  botUser.Username,
		TgFirstname: botUser.FirstName,
		TgLastname:  botUser.LastName,
	})
	return err
}

func (r *Users) Update(user *model.User) (err error) {
	ctx := context.Background()
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		status := "success"
		if err != nil {
			status = "failure"
		}
		attrs := []attribute.KeyValue{
			attribute.String("db.table", "users"),
			attribute.String("db.operation", "update_user_session"), // Assuming this is primarily for session updates
			attribute.String("status", status),
		}
		if dbOperationsCounter != nil {
			dbOperationsCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		if dbOperationDuration != nil {
			dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		}
	}()

	err = r.DB.SetUser(user)
	return err
}
