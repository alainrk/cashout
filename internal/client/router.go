package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"fmt"
	"strings"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
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

	if user.Session.State == model.StateInsertingIncome || user.Session.State == model.StateInsertingExpense {
		return c.addTransaction(b, ctx, user)
	}

	// During-insert edit transaction

	if user.Session.State == model.StateEditingTransactionDate {
		return c.editTransactionDate(b, ctx, user)
	}

	if user.Session.State == model.StateEditingTransactionCategory {
		return c.editTransactionCategory(b, ctx, user)
	}

	if user.Session.State == model.StateEditingTransactionAmount {
		return c.editTransactionAmount(b, ctx, user)
	}

	if user.Session.State == model.StateEditingTransactionDescription {
		return c.editTransactionDescription(b, ctx, user)
	}

	// End of during-insert edit transaction

	// Top-level edit transaction

	if user.Session.State == model.StateTopLevelEditingTransactionDate {
		return c.EditTransactionDateConfirm(b, ctx)
	}

	if user.Session.State == model.StateTopLevelEditingTransactionCategory {
		return c.EditTransactionCategoryConfirm(b, ctx)
	}

	if user.Session.State == model.StateTopLevelEditingTransactionAmount {
		return c.EditTransactionAmountConfirm(b, ctx)
	}

	if user.Session.State == model.StateTopLevelEditingTransactionDescription {
		return c.EditTransactionDescriptionConfirm(b, ctx)
	}

	// Search-related states
	if user.Session.State == model.StateEnteringSearchQuery {
		return c.SearchQueryEntered(b, ctx)
	}

	// End of top-level edit transaction

	// Default behavior: start transaction flow for any unhandled text.
	// Heuristic 1: there must be at least a digit in text.
	if strings.ContainsAny(ctx.Message.Text, "0123456789") {

		// Heuristic 2: it's more common to be an expense than an income, set it to default.
		user.Session.State = model.StateInsertingExpense
		// Heuristic 3: try to extract if it could be an income by looking in the text.
		if utils.IsAnIncomeTransactionPrompt(ctx.Message.Text) {
			user.Session.State = model.StateInsertingIncome
		}
		err = c.Repositories.Users.Update(&user)
		if err != nil {
			return fmt.Errorf("failed to set user data: %w", err)
		}

		return c.addTransaction(b, ctx, user)
	}

	c.CleanupKeyboard(b, ctx)
	c.SendHomeKeyboard(b, ctx, "Sorry I don't understand, what can I do for you?\n\n/edit - Edit a transaction\n/delete - Delete a transaction\n/search - Search transactions\n/list - List your transactions\n/week Week Recap\n/month Month Recap\n/year Year Recap\n/export - Export all transactions to CSV")

	return fmt.Errorf("invalid top-level state")
}
