package client

import (
	"errors"
	"fmt"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		_, errm := b.SendMessage(ctx.EffectiveChat.Id, fmt.Sprintf("You're not allowed to use this bot.\nContact the admin giving them your Telegram ID: \"%d\", and your username (you must create one): \"%s\"", ctx.EffectiveChat.Id, ctx.EffectiveChat.Username), nil)
		return errors.Join(err, errm)
	}

	msg := fmt.Sprintf("Welcome to Cashout, %s!\nWhat can I do for you?\n\n/list - List your transactions\n/edit - Edit a transaction\n/delete - Delete a transaction\n/search - Search transactions\n/week - Week Recap\n/month - Month Recap\n/year - Year Recap\n/export - Export all transactions to CSV\n/cancel - Cancel current operation\n/start - Show this menu again\n/new - Show this menu again", user.Name)

	err = c.SendHomeKeyboard(b, ctx, msg)

	return err
}
