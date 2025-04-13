package client

import (
	"fmt"
	"happypoor/internal/ai"
	"happypoor/internal/db"
	"happypoor/internal/model"
	"happypoor/internal/repository"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Client struct {
	Repositories Repositories
	LLM          ai.LLM
}

type Repositories struct {
	Users repository.Users
}

func NewClient(db *db.DB, llm ai.LLM) *Client {
	return &Client{
		Repositories: Repositories{
			Users: repository.Users{DB: db},
		},
		LLM: llm,
	}
}

// authAndGetUser authenticates the user and returns the user data.
func (c *Client) authAndGetUser(ctx *ext.Context) (model.User, error) {
	user, exists, err := c.Repositories.Users.GetByUsername(ctx.Message.From.Username)
	if err != nil {
		return user, fmt.Errorf("failed to get user data: %w", err)
	}

	if exists {
		user.Session.Iterations++
		c.Repositories.Users.Update(&user)
		return user, nil
	}

	// User to be created
	session := model.UserSession{
		Iterations:  0,
		State:       model.StateNormal,
		LastCommand: model.CommandStart,
		LastMessage: ctx.Message.Text,
	}

	if err := c.Repositories.Users.UpsertWithContext(ctx, session); err != nil {
		return user, fmt.Errorf("failed to set user data: %w", err)
	}

	user, _, err = c.Repositories.Users.GetByUsername(ctx.Message.From.Username)
	if err != nil {
		return user, fmt.Errorf("failed to get user data: %w", err)
	}

	return user, nil
}
