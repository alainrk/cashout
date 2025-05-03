package client

import (
	"cashout/internal/model"
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func (c *Client) FreeTextRouter(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	msg := ctx.Message
	if message.Text(msg) && strings.ToLower(strings.Trim(msg.Text, " ")) == "cancel" {
		return c.Cancel(b, ctx)
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	if user.Session.State == model.StateInsertingIncome || user.Session.State == model.StateInsertingExpense {
		return c.addTransaction(b, ctx, user)
	}

	if user.Session.State == model.StateEditingTransactionDate {
		return c.editTransactionDate(b, ctx, user)
	}

	if user.Session.State == model.StateEditingTransactionCategory {
		return c.editTransactionCategory(b, ctx, user)
	}

	c.CleanupKeyboard(b, ctx)
	c.SendHomeKeyboard(b, ctx, "Sorry I don't understand, what can I do for you?\n\n/delete - Delete a transaction\n/list - List your transactions\n/month Month Recap\n/year Year Recap")

	return fmt.Errorf("invalid top-level state")
}
