package client

import (
	"fmt"
	"happypoor/internal/model"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Welcome to HappyPoor, %s!", user.Name)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}

// Message handles incoming messages not in a specific flow.
func (c *Client) Message(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	user.Session.State = model.StateWaiting
	// TODO: This must be understood through the analysis of the text the user have given
	user.Session.LastCommand = model.CommandAddIncome
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	// TODO: Change type here according to buttons, previous iteration with the bot
	transaction, err := c.LLM.ExtractTransaction(ctx.Message.Text, model.TypeIncome)
	if err != nil {
		msg := fmt.Sprintf("I'm sorry, I can't understand your transaction '%s', %s!", user.Session.LastMessage, user.Name)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return err
	}

	msg := fmt.Sprintf("%s (€ %.2f), %s", transaction.Category, transaction.Amount, transaction.Description)
	fmt.Println(msg)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastCommand = model.CommandCancel
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg := fmt.Sprintf("Sure, you can always restart with /start, %s!", user.Name)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
