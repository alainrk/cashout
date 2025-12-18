package client

import (
	"errors"
	"fmt"
	"strings"

	"cashout/internal/ai"
	"cashout/internal/model"
	"cashout/internal/utils"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// FreeTextRouter is the default router for the bot.
// It tries to infer if the user is adding a transaction (and in that case if it is an expense or income) if the state is not set to any other state.
// If the user is adding a transaction, it sets the correct prompt and calls the LLM to extract the transaction information.
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

	// The use pre-selected the adding transaction flow, so we don't need to infer it.
	if user.Session.State == model.StateInsertingIncome || user.Session.State == model.StateInsertingExpense {
		return c.addTransaction(b, ctx, user)
	}

	// During-insert edit transaction.

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

	// Top-level edit transaction cases.

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

	// Search-related states.
	if user.Session.State == model.StateEnteringSearchQuery {
		return c.SearchQueryEntered(b, ctx)
	}

	// Edit search-related states.
	if user.Session.State == model.StateEnteringEditSearchQuery {
		return c.EditSearchQueryEntered(b, ctx)
	}

	// Delete search-related states.
	if user.Session.State == model.StateEnteringDeleteSearchQuery {
		return c.DeleteSearchQueryEntered(b, ctx)
	}

	// Free text top level case: use LLM to classify user intent.
	return c.classifyAndRouteIntent(b, ctx, user)
}

// classifyAndRouteIntent uses the LLM to classify the user's intent and routes to the appropriate handler
func (c *Client) classifyAndRouteIntent(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	// Quick heuristic: if text contains digits, it's likely a transaction
	// Use fast local check before calling LLM.
	if strings.ContainsAny(ctx.Message.Text, "0123456789") {

		// Default to expense and check for income keywords.
		user.Session.State = model.StateInsertingExpense
		if utils.IsAnIncomeTransactionPrompt(ctx.Message.Text) {
			user.Session.State = model.StateInsertingIncome
		}

		err := c.Repositories.Users.Update(&user)
		if err != nil {
			return fmt.Errorf("failed to set user data: %w", err)
		}

		return c.addTransaction(b, ctx, user)
	}

	// Call LLM to classify intent for any other case.
	classifiedIntent, err := c.LLM.ClassifyIntent(ctx.Message.Text)
	if err != nil {
		c.Logger.Warnf("Failed to classify intent: %v, falling back to unknown", err)
		classifiedIntent = ai.ClassifiedIntent{Intent: ai.IntentUnknown, Confidence: 0}
	}

	c.Logger.Debugf("Classified intent: %s (confidence: %.2f)", classifiedIntent.Intent, classifiedIntent.Confidence)

	// Route based on classified intent
	switch classifiedIntent.Intent {
	case ai.IntentAddExpense:
		return c.AddTransactionExpense(b, ctx)

	case ai.IntentAddIncome:
		return c.AddTransactionIncome(b, ctx)

	case ai.IntentEdit:
		return c.EditTransactions(b, ctx)

	case ai.IntentDelete:
		return c.DeleteTransactions(b, ctx)

	case ai.IntentSearch:
		return c.SearchTransactions(b, ctx)

	case ai.IntentList:
		return c.ListTransactions(b, ctx)

	case ai.IntentWeekRecap:
		return c.WeekRecap(b, ctx)

	case ai.IntentMonthRecap:
		return c.MonthRecap(b, ctx)

	case ai.IntentYearRecap:
		return c.YearRecap(b, ctx)

	case ai.IntentExport:
		return c.ExportTransactions(b, ctx)

	default:
		// Unknown intent - show help
		err = c.CleanupKeyboard(b, ctx)
		err = errors.Join(err, c.SendHomeKeyboard(b, ctx, "I'm not sure what you'd like to do. Here are the available options:\n\n/edit - Edit a transaction\n/delete - Delete a transaction\n/search - Search transactions\n/list - List your transactions\n/week - Week Recap\n/month - Month Recap\n/year - Year Recap\n/export - Export all transactions to CSV\n\nOr just type a transaction like \"coffee 5\" to add it!"))
		if err != nil {
			return err
		}
		return nil
	}
}
