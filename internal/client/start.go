package client

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Welcome to Cashout, %s!\nWhat can I do for you?\n\n/edit - Edit a transaction\n/delete - Delete a transaction\n/list - List your transactions\n/week Week Recap\n/month Month Recap\n/year Year Recap\n/export - Export all transactions to CSV", user.Name)

	c.SendHomeKeyboard(b, ctx, msg)

	return nil
}
