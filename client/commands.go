package client

import (
	"fmt"
	"happypoor/internal/model"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// AuthAndGetUser authenticates the user and returns the user data.
func (c *Client) AuthAndGetUser(b *gotgbot.Bot, ctx *ext.Context) (model.User, error) {
	user, exists, err := c.getUserData(ctx, ctx.Message.From.Username)
	if err != nil {
		return user, fmt.Errorf("failed to get user data: %w", err)
	}

	if exists {
		user.Session.Iterations++
		c.Db.SetUser(&user)
		return user, nil
	}

	// User to be created
	session := model.UserSession{
		Iterations:  0,
		State:       model.StateNormal,
		LastCommand: model.CommandStart,
		LastMessage: ctx.Message.Text,
	}

	if err := c.setUserData(ctx, session); err != nil {
		return user, fmt.Errorf("failed to set user data: %w", err)
	}

	user, _, err = c.getUserData(ctx, ctx.Message.From.Username)
	if err != nil {
		return user, fmt.Errorf("failed to get user data: %w", err)
	}

	return user, nil
}

// Start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.AuthAndGetUser(b, ctx)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Welcome to HappyPoor %s!", user.Name)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return handlers.NextConversationState("normal")
}

// Message handles incoming messages not in a specific flow.
func (c *Client) Message(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.AuthAndGetUser(b, ctx)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	return nil
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.AuthAndGetUser(b, ctx)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastCommand = model.CommandCancel
	user.Session.LastMessage = ctx.Message.Text

	err = c.Db.SetUser(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	return nil
}
