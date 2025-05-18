package client

import (
	"cashout/internal/model"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// authAndGetUser authenticates the user and returns the user data.
func (c *Client) authAndGetUser(user gotgbot.User) (model.User, error) {
	if c.Config.AuthEnabled {
		if _, ok := c.Config.AllowedUsers[user.Username]; !ok {
			return model.User{}, fmt.Errorf("user %s is not allowed", user.Username)
		}
	}

	u, exists, err := c.Repositories.Users.GetByUsername(user.Username)
	if err != nil {
		return u, fmt.Errorf("failed to get user data: %w", err)
	}

	if exists {
		c.Repositories.Users.Update(&u)
		return u, nil
	}

	// First Message, user to be created.
	session := model.UserSession{
		State: model.StateStart,
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

func (c *Client) getUserFromContext(ctx *ext.Context) (isInline bool, user gotgbot.User) {
	if ctx.CallbackQuery != nil {
		return true, ctx.CallbackQuery.From
	}
	return false, *ctx.Message.From
}
