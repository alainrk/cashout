package client

import (
	"encoding/json"
	"fmt"
	"happypoor/internal/model"
	"happypoor/internal/utils"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) AddTransactionIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	query := ctx.CallbackQuery
	msg := query.Message

	action := strings.Split(query.Data, ".")[2]

	switch action {
	case "income":
		user.Session.State = model.StateInsertingIncome
	case "expense":
		user.Session.State = model.StateInsertingExpense
	default:
		c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
		return nil
	}

	user.Session.LastMessage = action

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg.EditText(b, fmt.Sprintf("Sure, to add a new <b>%s</b>:\nTell me category, amount and description. You can also specify a date and change it later, today is default.\n\n<i>Examples:</i>\n<code>Irish Pub 3.4</code>\n<code>January salary 3k 10/01</code>", action), &gotgbot.EditMessageTextOpts{
		ParseMode: "HTML",
	})
	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Cancel", CallbackData: "transactions.cancel"}},
			},
		},
	})

	return ext.ContinueGroups
}

func (c *Client) addTransaction(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transactionType model.TransactionType

	switch user.Session.State {
	case model.StateInsertingIncome:
		transactionType = model.TypeIncome
	case model.StateInsertingExpense:
		transactionType = model.TypeExpense
	default:
		c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
		return nil
	}

	transaction, err := c.LLM.ExtractTransaction(ctx.Message.Text, transactionType)
	if err != nil {
		msg := fmt.Sprintf("I'm sorry, I can't understand your transaction '%s', %s!", user.Session.LastMessage, user.Name)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return err
	}

	if transaction.Amount == 0 {
		msg := fmt.Sprintf("I'm sorry, I can't understand your transaction '%s', %s!", user.Session.LastMessage, user.Name)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Store the transaction in the session
	user.Session.State = model.StateNormal
	user.Session.LastMessage = ctx.Message.Text
	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, msg, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("failed to send confirm message: %v\n", err)
		return err
	}

	return nil
}

func (c *Client) EditTransactionIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	c.CleanupInlineKeyboard(b, ctx)

	user.Session.State = model.StateEditingTransaction
	user.Session.LastMessage = "edit"

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	b.SendMessage(ctx.EffectiveSender.ChatId, "Add your date (e.g. dd mm, dd-mm, dd-mm-yyyy).", &gotgbot.SendMessageOpts{})

	return nil
}

// EditTransaction edits a transaction.
func (c *Client) editTransaction(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	// TODO: Handle this with LLM to support different formats automatically and also yesterday, 2 days ago etc.

	// Get date from DD-MM-YYYY to date
	date, err := utils.ParseDate(ctx.Message.Text)
	if err != nil {
		// TODO: Handle invalid date
		fmt.Printf("failed to parse date: %v\n", err)
		return err
	}

	// TODO: Reject future date

	transaction.Date = date

	user.Session.LastMessage = ctx.Message.Text
	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	m := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
				},
			},
		},
	})

	return nil
}

// Confirm confirms the previous action after the user been prompted.
func (c *Client) Confirm(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	var transaction model.Transaction
	err = json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	user.Session.State = model.StateNormal
	user.Session.LastMessage = "confirm"
	user.Session.Body = ""

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	transaction.TgID = user.TgID
	transaction.Currency = model.CurrencyEUR

	err = c.Repositories.Transactions.Add(transaction)
	if err != nil {
		// TODO: Handle send failure message to the user
		fmt.Printf("failed to add transaction: %v\n", err)
		return fmt.Errorf("failed to add transaction: %w", err)
	}

	// Remove the keyboard from the previous message
	query := ctx.CallbackQuery
	msg := query.Message
	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{},
		},
	})

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}
	c.SendAddTransactionKeyboard(b, ctx, fmt.Sprintf("%s Your transaction has been saved!", emoji))

	return nil
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery
	msg := query.Message

	user.Session.State = model.StateNormal
	user.Session.LastMessage = "cancel"

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
	})

	c.SendAddTransactionKeyboard(b, ctx, "Add a transaction")

	return ext.EndGroups
}
