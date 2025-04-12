package client

import (
	"fmt"
	"happypoor/internal/db"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	// Get the existing data. "ok" will be false if the data doesn't exist yet.
	user, exists, err := c.getUserData(ctx, ctx.Message.From.Username)
	if err != nil {
		return fmt.Errorf("failed to get user data: %w", err)
	}

	if exists {
		msg := fmt.Sprintf("Your id is <code>%d</code>", user.TgID)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	session := db.UserSession{
		State:       db.StateNormal,
		LastCommand: db.CommandStart,
		LastMessage: ctx.Message.Text,
	}

	if err := c.setUserData(ctx, session); err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	user, _, err = c.getUserData(ctx, ctx.Message.From.Username)
	if err != nil {
		return fmt.Errorf("failed to get user data: %w", err)
	}

	msg := fmt.Sprintf("Welcome to Happypoor %s!", user.Name)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
