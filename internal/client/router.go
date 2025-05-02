package client

import (
	"cashout/internal/model"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) FreeTextRouter(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
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

	return fmt.Errorf("invalid top-level state")
}
