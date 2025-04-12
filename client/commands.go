package client

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	// Get the existing data. "ok" will be false if the data doesn't exist yet.
	countVal, ok := c.getUserData(ctx, "count")
	if !ok {
		c.setUserData(ctx, "count", 1)
		ctx.EffectiveMessage.Reply(b, "This is the first time you press start.", &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Cast the data to an int so it can be used.
	count, ok := countVal.(int)
	if !ok {
		ctx.EffectiveMessage.Reply(b, "'count' was not an integer, as was expected - this is a programmer error!", &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Increment our count (one more press!)
	count += 1
	c.setUserData(ctx, "count", count)
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("You have pressed start %d times.", count), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	// Collect user information
	msg := fmt.Sprintf("%s, %d, %s %s", ctx.Message.From.Username, ctx.Message.From.Id, ctx.Message.From.FirstName, ctx.Message.From.LastName)
	_, err := ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}
