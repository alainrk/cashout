package client

import (
	"cashout/internal/model"
	"context"
)

type contextKey string

const userContextKey contextKey = "user"

// SetUserInContext adds a user to the context
func SetUserInContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUserFromContext retrieves a user from the context
func GetUserFromContext(ctx context.Context) *model.User {
	user, ok := ctx.Value(userContextKey).(*model.User)
	if !ok {
		return nil
	}
	return user
}
