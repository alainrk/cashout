package client

import (
	"fmt"
	"happypoor/internal/ai"
	"happypoor/internal/db"
	"happypoor/internal/model"
	"happypoor/internal/repository"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Client struct {
	Repositories Repositories
	LLM          ai.LLM
}

type Repositories struct {
	Users        repository.Users
	Transactions repository.Transactions
}

func NewClient(db *db.DB, llm ai.LLM) *Client {
	return &Client{
		Repositories: Repositories{
			Users: repository.Users{DB: db},
			Transactions: repository.Transactions{
				DB: db,
			},
		},
		LLM: llm,
	}
}

// authAndGetUser authenticates the user and returns the user data.
func (c *Client) authAndGetUser(user gotgbot.User) (model.User, error) {
	u, exists, err := c.Repositories.Users.GetByUsername(user.Username)
	if err != nil {
		return u, fmt.Errorf("failed to get user data: %w", err)
	}

	if exists {
		u.Session.Iterations++
		c.Repositories.Users.Update(&u)
		return u, nil
	}

	// First Message, user to be created.
	session := model.UserSession{
		Iterations: 0,
		State:      model.StateStart,
	}

	if err := c.Repositories.Users.UpsertWithContext(user, session); err != nil {
		return u, fmt.Errorf("failed to set user data: %w", err)
	}

	u, _, err = c.Repositories.Users.GetByUsername(user.Username)
	if err != nil {
		return u, fmt.Errorf("failed to get user data: %w", err)
	}

	return u, nil
}

func (c *Client) getUserFromContext(ctx *ext.Context) (isTopLevel bool, user gotgbot.User) {
	if ctx.CallbackQuery != nil {
		return false, ctx.CallbackQuery.From
	}
	return true, *ctx.Message.From
}
